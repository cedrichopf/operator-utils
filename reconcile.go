package utils

import (
	"context"
	"time"

	"github.com/cedrichopf/operator-utils/hash"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type reconcileParams struct {
	context  context.Context
	client   client.Client
	scheme   *runtime.Scheme
	kind     string
	expected client.Object
	owner    metav1.Object
	// Temporary object for fetching object state from cluster
	tmpObj client.Object
}

type ReconcileResult struct {
	Updated      bool
	Requeue      bool
	RequeueAfter time.Duration
}

func reconcileResource(params *reconcileParams) (ReconcileResult, error) {
	log := log.FromContext(params.context)

	kind := params.kind
	name := params.expected.GetName()
	namespace := params.expected.GetNamespace()

	revisionHash, err := hash.GenerateObjectHash(params.expected)
	if err != nil {
		log.Error(
			err, "Unable to generate revision hash for object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)
		return ReconcileResult{}, err
	}

	err = params.client.Get(params.context, types.NamespacedName{Name: name, Namespace: namespace}, params.tmpObj)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)

		// Add revision hash label for reconcile
		labels := params.expected.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[hash.REVISION_HASH_LABEL] = revisionHash
		params.expected.SetLabels(labels)

		err = ctrl.SetControllerReference(params.owner, params.expected, params.scheme)
		if err != nil {
			return ReconcileResult{}, err
		}

		err = params.client.Create(params.context, params.expected)
		if err != nil {
			log.Error(
				err,
				"Failed to create new object",
				"Object.Kind", kind,
				"Object.Name", name,
				"Object.Namespace", namespace,
			)
			return ReconcileResult{
				Updated:      false,
				Requeue:      true,
				RequeueAfter: 10 * time.Second,
			}, err
		}
		return ReconcileResult{
			Updated: true,
		}, nil
	} else if err != nil {
		log.Error(
			err,
			"Failed to get object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)
		return ReconcileResult{}, err
	}

	// Check object
	labels := params.tmpObj.GetLabels()
	if revisionHash != labels[hash.REVISION_HASH_LABEL] {
		log.Info(
			"Updating outdated object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)

		labels = params.expected.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[hash.REVISION_HASH_LABEL] = revisionHash
		params.expected.SetLabels(labels)

		err = params.client.Update(params.context, params.expected)
		if err != nil {
			log.Error(
				err,
				"Unable to update object",
				"Object.Kind", kind,
				"Object.Name", name,
				"Object.Namespace", namespace,
			)
			return ReconcileResult{}, err
		}
		return ReconcileResult{
			Updated: true,
		}, nil
	}
	return ReconcileResult{}, nil
}
