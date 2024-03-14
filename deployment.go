package utils

import (
	"context"
	"maps"
	"time"

	"github.com/cedrichopf/operator-utils/hash"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ReconcileDeployment(ctx context.Context, expected *appsv1.Deployment, owner metav1.Object, c client.Client, s *runtime.Scheme) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	revisionHash, err := hash.GenerateObjectHash(expected)
	if err != nil {
		log.Error(
			err, "Unable to generate revision hash for deployment",
			"Deployment.Name", expected.Name,
			"Deployment.Namespace", expected.Namespace,
		)
		return nil, err
	}
	hashLabel := hash.GenerateRevisionHashLabel(revisionHash)

	deployment := &appsv1.Deployment{}
	err = c.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, deployment)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new deployment",
			"Deployment.Name", expected.Name,
			"Deployment.Namespace", expected.Namespace,
		)

		// Add revision hash label for reconcile
		maps.Copy(expected.Labels, hashLabel)

		err = ctrl.SetControllerReference(owner, expected, s)
		if err != nil {
			return nil, err
		}

		if err = c.Create(ctx, expected); err != nil {
			log.Error(
				err,
				"Failed to create new deployment",
				"Deployment.Name", expected.Name,
				"Deployment.Namespace", expected.Namespace,
			)

			return &ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}

		return nil, nil
	} else if err != nil {
		log.Error(
			err,
			"Failed to get deployment",
			"Deployment.Name", expected.Name,
			"Deployment.Namespace", expected.Namespace,
		)
		return nil, err
	}

	// Check deployment
	if revisionHash != deployment.Labels[hash.REVISION_HASH_LABEL] {
		log.Info(
			"Updating outdated deployment",
			"Deployment.Name", expected.Name,
			"Deployment.Namespace", expected.Namespace,
		)

		maps.Copy(expected.Labels, hashLabel)

		err = c.Update(ctx, expected)
		if err != nil {
			log.Error(
				err,
				"Unable to update deployment",
				"Deployment.Name", expected.Name,
				"Deployment.Namespace", expected.Namespace,
			)
			return nil, err
		}
	}

	return nil, nil
}
