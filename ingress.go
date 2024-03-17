package utils

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileIngress(ctx context.Context, expected *networkingv1.Ingress, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&reconcileParams{
		context:  ctx,
		client:   client,
		scheme:   scheme,
		kind:     "Ingress",
		expected: expected,
		owner:    owner,
		tmpObj:   &networkingv1.Ingress{},
	})
}
