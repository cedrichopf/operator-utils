package utils

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileDeployment takes a Deployment and creates or updates it on a Kubernetes cluster. It will generate
// a revision hash that will be added to the object as a label and add the controller reference to the given
// owner object.
//
// A successful ReconcileDeployment returns a ReconcileResult and nil.
// A failed ReconcileDeployment returns an empty ReconsileResult and an error.
func ReconcileDeployment(ctx context.Context, expected *appsv1.Deployment, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&ReconcileParams{
		Context:  ctx,
		Client:   client,
		Scheme:   scheme,
		kind:     "Deployment",
		Expected: expected,
		Owner:    owner,
		tmpObj:   &appsv1.Deployment{},
	})
}
