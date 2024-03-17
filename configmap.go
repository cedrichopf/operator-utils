package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileConfigMap(ctx context.Context, expected *corev1.ConfigMap, owner metav1.Object, c client.Client, s *runtime.Scheme) (ReconcileResult, error) {
	return reconcileResource(&ReconcileParams{
		Context:  ctx,
		Client:   c,
		Scheme:   s,
		Kind:     "ConfigMap",
		Expected: expected,
		tmpObj:   &corev1.ConfigMap{},
		Owner:    owner,
	})
}
