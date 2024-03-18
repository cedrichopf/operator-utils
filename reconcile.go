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

// ReconcileParams contains all required information to reconcile an object
type ReconcileParams struct {
	Context  context.Context
	Client   client.Client
	Scheme   *runtime.Scheme
	Expected client.Object
	Owner    metav1.Object

	kind string
	// Temporary object for fetching object state from cluster
	tmpObj client.Object
}

// ReconcileResult contains information about a reconciled resource.
//
// Updated shows that the reconciled resource has been updated on the Kubernetes cluster.
// Requeue shows that the operator should requeue the current reconcile.
// RequeueAfter provides a timeout which should be used by the operator to requeue the current reconcile.
type ReconcileResult struct {
	Updated      bool
	Requeue      bool
	RequeueAfter time.Duration
}

// reconcileResource is a generic reconcile function that can be used by several object kinds.
//
// A successful reconcileResource returns a ReconcileResult and nil.
// A failed reconcileResource returns an empty ReconsileResult and an error.
func reconcileResource(params *ReconcileParams) (ReconcileResult, error) {
	log := log.FromContext(params.Context)

	kind := params.kind
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
