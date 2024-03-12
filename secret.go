package utils

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ReconcileSecret(ctx context.Context, expected *corev1.Secret, owner metav1.Object, c client.Client, s *runtime.Scheme) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	secret := &corev1.Secret{}
	err := c.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, secret)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new secret",
			"Secret.Name", expected.Name,
			"Secret.Namespace", expected.Namespace,
		)

		err = ctrl.SetControllerReference(owner, expected, s)
		if err != nil {
			return nil, err
		}

		if err = c.Create(ctx, expected); err != nil {
			log.Error(
				err,
				"Failed to create new secret",
				"Secret.Name", expected.Name,
				"Secret.Namespace", expected.Namespace,
			)

			return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}
	} else if err != nil {
		log.Error(
			err,
			"Failed to get secret",
			"Secret.Name", expected.Name,
			"Secret.Namespace", expected.Namespace,
		)
		return nil, err
	}

	//TODO check state

	return nil, nil
}
