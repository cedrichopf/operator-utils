package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileSecret(ctx context.Context, expected *corev1.Secret, owner metav1.Object, client client.Client, scheme *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&reconcileParams{
		context:  ctx,
		client:   client,
		scheme:   scheme,
		kind:     "Secret",
		expected: expected,
		owner:    owner,
		tmpObj:   &corev1.Secret{},
	})
}
