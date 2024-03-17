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

type ReconcileParams struct {
	Context  context.Context
	Client   client.Client
	Scheme   *runtime.Scheme
	Kind     string
	Expected client.Object
	Owner    metav1.Object

	// Temporary object for fetching object state from cluster
	tmpObj client.Object
}

type ReconcileResult struct {
	Updated      bool
	Requeue      bool
	RequeueAfter time.Duration
}

func reconcileResource(params *ReconcileParams) (ReconcileResult, error) {
	log := log.FromContext(params.Context)

	kind := params.Kind
	name := params.Expected.GetName()
	namespace := params.Expected.GetNamespace()

	revisionHash, err := hash.GenerateObjectHash(params.Expected)
	if err != nil {
		log.Error(
			err, "Unable to generate revision hash for object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)
		return ReconcileResult{}, err
	}

	err = params.Client.Get(params.Context, types.NamespacedName{Name: name, Namespace: namespace}, params.tmpObj)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new object",
			"Object.Kind", kind,
			"Object.Name", name,
			"Object.Namespace", namespace,
		)

		// Add revision hash label for reconcile
		labels := params.Expected.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[hash.REVISION_HASH_LABEL] = revisionHash
		params.Expected.SetLabels(labels)

		err = ctrl.SetControllerReference(params.Owner, params.Expected, params.Scheme)
		if err != nil {
			return ReconcileResult{}, err
		}

		err = params.Client.Create(params.Context, params.Expected)
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

		labels = params.Expected.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[hash.REVISION_HASH_LABEL] = revisionHash
		params.Expected.SetLabels(labels)

		err = params.Client.Update(params.Context, params.Expected)
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
