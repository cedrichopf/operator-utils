package utils

import (
	"context"
	"maps"
	"time"

	"github.com/cedrichopf/operator-utils/hash"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ReconcileConfigMap(ctx context.Context, expected *corev1.ConfigMap, owner metav1.Object, c client.Client, s *runtime.Scheme) (*ctrl.Result, error) {
	log := log.FromContext(ctx)

	revisionHash, err := hash.GenerateObjectHash(expected)
	if err != nil {
		log.Error(
			err, "Unable to generate revision hash for configmap",
			"ConfigMap.Name", expected.Name,
			"ConfigMap.Namespace", expected.Namespace,
		)
		return nil, err
	}
	hashLabel := hash.GenerateRevisionHashLabel(revisionHash)

	configMap := &corev1.ConfigMap{}
	err = c.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, configMap)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info(
			"Creating new configmap",
			"ConfigMap.Name", expected.Name,
			"ConfigMap.Namespace", expected.Namespace,
		)

		// Add revision hash label for reconcile
		if expected.Labels == nil {
			expected.Labels = make(map[string]string)
		}
		maps.Copy(expected.Labels, hashLabel)

		err = ctrl.SetControllerReference(owner, expected, s)
		if err != nil {
			return nil, err
		}

		if err = c.Create(ctx, expected); err != nil {
			log.Error(
				err,
				"Failed to create new configmap",
				"ConfigMap.Name", expected.Name,
				"ConfigMap.Namespace", expected.Namespace,
			)

			return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}
		return nil, nil
	} else if err != nil {
		log.Error(
			err,
			"Failed to get configmap",
			"ConfigMap.Name", expected.Name,
			"ConfigMap.Namespace", expected.Namespace,
		)
		return nil, err
	}

	// Check configmap
	if revisionHash != configMap.Labels[hash.REVISION_HASH_LABEL] {
		log.Info(
			"Updating outdated configmap",
			"ConfigMap.Name", expected.Name,
			"ConfigMap.Namespace", expected.Namespace,
		)

		if expected.Labels == nil {
			expected.Labels = make(map[string]string)
		}
		maps.Copy(expected.Labels, hashLabel)

		err = c.Update(ctx, expected)
		if err != nil {
			log.Error(
				err,
				"Unable to update configmap",
				"ConfigMap.Name", expected.Name,
				"ConfigMap.Namespace", expected.Namespace,
			)
			return nil, err
		}
	}

	return nil, nil
}