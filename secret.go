package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileSecret takes a Secret and creates or updates it on a Kubernetes cluster. It will generate a
// revision hash that will be added to the object as a label and add the controller reference to the given
// owner object.
//
// A successful ReconcileSecret returns a ReconcileResult and nil.
// A failed ReconcileSecret returns an empty ReconsileResult and an error.
func ReconcileSecret(ctx context.Context, expected *corev1.Secret, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&ReconcileParams{
		Context:  ctx,
		Client:   client,
		Scheme:   scheme,
		kind:     "Secret",
		Expected: expected,
		Owner:    owner,
		tmpObj:   &corev1.Secret{},
	})
}
