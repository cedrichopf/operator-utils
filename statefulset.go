package utils

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ReconcileStatefulSet(ctx context.Context, expected *appsv1.StatefulSet, owner metav1.Object, c client.Client, s *runtime.Scheme) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	statefulset := &appsv1.StatefulSet{}
	err := c.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, statefulset)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new statefulset",
			"StatefulSet.Name", expected.Name,
			"StatefulSet.Namespace", expected.Namespace,
		)

		err = ctrl.SetControllerReference(owner, expected, s)
		if err != nil {
			return nil, err
		}

		if err = c.Create(ctx, expected); err != nil {
			log.Error(
				err,
				"Failed to create new statefulset",
				"StatefulSet.Name", expected.Name,
				"StatefulSet.Namespace", expected.Namespace,
			)

			return &ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		log.Error(
			err,
			"Failed to get statefulset",
			"Deployment.Name", expected.Name,
			"Deployment.Namespace", expected.Namespace,
		)
		return nil, err
	}

	//TODO check state

	return nil, nil
}
