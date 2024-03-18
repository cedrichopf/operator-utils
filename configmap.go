package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileConfigMap takes a ConfigMap and creates or updates it on a Kubernetes cluster. It will generate a
// revision hash that will be added to the object as a label and add the controller reference to the given
// owner object.
//
// A successful ReconcileConfigMap returns a ReconcileResult and nil.
// A failed ReconcileConfigMap returns an empty ReconsileResult and an error.
func ReconcileConfigMap(ctx context.Context, expected *corev1.ConfigMap, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&ReconcileParams{
		Context:  ctx,
		Client:   client,
		Scheme:   scheme,
		kind:     "ConfigMap",
		Expected: expected,
		Owner:    owner,
		tmpObj:   &corev1.ConfigMap{},
	})
}
