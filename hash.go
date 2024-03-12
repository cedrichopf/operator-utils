package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"maps"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
		"app.kubernetes.io/revision-hash": revisionHash,
	}
}

func AddRevisionHashLabel(object metav1.Object) error {
	revisionHash, err := GenerateObjectHash(object)
	if err != nil {
		return err
	}
	hashLabel := GenerateRevisionHashLabel(revisionHash)
	labels := object.GetLabels()
	maps.Copy(labels, hashLabel)
	object.SetLabels(labels)
	return nil
}
