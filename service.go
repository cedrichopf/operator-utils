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

func ReconcileService(ctx context.Context, expected *corev1.Service, owner metav1.Object, c client.Client, s *runtime.Scheme) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	service := &corev1.Service{}
	err := c.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, service)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new service",
			"Service.Name", expected.Name,
			"Service.Namespace", expected.Namespace,
		)

		err = ctrl.SetControllerReference(owner, expected, s)
		if err != nil {
			return nil, err
		}

		if err = c.Create(ctx, expected); err != nil {
			log.Error(
				err,
				"Failed to create new service",
				"Service.Name", expected.Name,
				"Service.Namespace", expected.Namespace,
			)

			return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}
	} else if err != nil {
		log.Error(
			err,
			"Failed to get service",
			"Service.Name", expected.Name,
			"Service.Namespace", expected.Namespace,
		)
		return nil, err
	}

	//TODO check state

	return nil, nil
}
