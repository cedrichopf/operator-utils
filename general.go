package utils

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteResources(ctx context.Context, c client.Client, objects []client.Object) error {
	for _, object := range objects {
		err := c.Delete(ctx, object)
		if err != nil && apierrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
	}
	return nil
}

func RemoveFinalizer(object metav1.Object, finalizer string) error {
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
		return fmt.Errorf("finalizer %s not found on object %s", finalizer, object.GetName())
	}

	object.SetFinalizers(labels)
	return nil
}
