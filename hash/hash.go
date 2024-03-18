package hash

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const REVISION_HASH_LABEL = "app.kubernetes.io/revision-hash"

// GenerateObjectHash generates a hash value for a given object
func GenerateObjectHash(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("unable to marshal object: %v", err)
	}

	hash := fnv.New32a()
	hash.Write(data)
	return fmt.Sprintf("%d", hash.Sum32()), nil
}

// GenerateRevisionHashLabel generates a map containing a given revision hash
func GenerateRevisionHashLabel(revisionHash string) map[string]string {
	return map[string]string{
		REVISION_HASH_LABEL: revisionHash,
	}
}

// AddRevisionHashLabel generates a revision hash and updates the labels of a given object
func AddRevisionHashLabel(object metav1.Object) error {
	revisionHash, err := GenerateObjectHash(object)
	if err != nil {
		return err
	}
	hashLabel := GenerateRevisionHashLabel(revisionHash)
	labels := object.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	for key, value := range hashLabel {
		labels[key] = value
	}
	object.SetLabels(labels)
	return nil
}
