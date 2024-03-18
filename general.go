package utils

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeleteResources takes a slice of Kubernetes objects and deletes them from a Kubernetes cluster
//
// DeleteResources returns an error if it fails to delete the object
func DeleteResources(ctx context.Context, client client.Client, objects []client.Object) error {
	for _, object := range objects {
		err := client.Delete(ctx, object)
		if err != nil && apierrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFinalizer deletes a given finalizer from a Kubernetes object
//
// RemoveFinalizer returns a bool value that indicates if the given finalizer has been removed
func RemoveFinalizer(object metav1.Object, finalizer string) bool {
	var found bool
	var labels []string
	for _, label := range object.GetFinalizers() {
		if label == finalizer {
			found = true
			continue
		}
		labels = append(labels, label)
	}

	if !found {
		return false
	}

	object.SetFinalizers(labels)
	return true
}
