package utils

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileIngress takes an Ingress and creates or updates it on a Kubernetes cluster. It will generate
// a revision hash that will be added to the object as a label and add the controller reference to the given
// owner object.
//
// A successful ReconcileIngress returns a ReconcileResult and nil.
// A failed ReconcileIngress returns an empty ReconsileResult and an error.
func ReconcileIngress(ctx context.Context, expected *networkingv1.Ingress, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&ReconcileParams{
		Context:  ctx,
		Client:   client,
		Scheme:   scheme,
		kind:     "Ingress",
		Expected: expected,
		Owner:    owner,
		tmpObj:   &networkingv1.Ingress{},
	})
}
