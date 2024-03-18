package utils

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileStatefulSet(ctx context.Context, expected *appsv1.StatefulSet, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&reconcileParams{
		context:  ctx,
		client:   client,
		scheme:   scheme,
		kind:     "StatefulSet",
		expected: expected,
		owner:    owner,
		tmpObj:   &appsv1.StatefulSet{},
	})
}
