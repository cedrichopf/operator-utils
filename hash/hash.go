package hash

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const REVISION_HASH_LABEL = "app.kubernetes.io/revision-hash"

func GenerateObjectHash(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("unable to marshal object: %v", err)
	}

	hash := fnv.New32a()
	hash.Write(data)
	return fmt.Sprintf("%d", hash.Sum32()), nil
}

func GenerateRevisionHashLabel(revisionHash string) map[string]string {
	return map[string]string{
		REVISION_HASH_LABEL: revisionHash,
	}
}

func AddRevisionHashLabel(object metav1.Object) error {
	revisionHash, err := GenerateObjectHash(object)
	if err != nil {
		return err
	}
	hashLabel := GenerateRevisionHashLabel(revisionHash)
	labels := object.GetLabels()
	for key, value := range hashLabel {
		labels[key] = value
	}
	object.SetLabels(labels)
	return nil
}
