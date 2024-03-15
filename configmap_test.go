package utils_test

import (
	"context"

	utils "github.com/cedrichopf/operator-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configMapName = "sample-configmap"
)

var _ = Describe("ConfigMap", func() {
	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &corev1.ConfigMap{}, client.InNamespace("default"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("creates a new ConfigMap", func() {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: namespace,
			},
		}

		result, err := utils.ReconcileConfigMap(context.Background(), configMap, owner, k8sClient, testEnv.Scheme)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("updates an outdated ConfigMap", func() {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: namespace,
			},
			Data: map[string]string{
				"sample": "data",
			},
		}

		result, err := utils.ReconcileConfigMap(context.Background(), configMap, owner, k8sClient, testEnv.Scheme)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeNil())

		updatedConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: namespace,
			},
			Data: map[string]string{
				"sample": "newData",
			},
		}

		result, err = utils.ReconcileConfigMap(context.Background(), updatedConfigMap, owner, k8sClient, testEnv.Scheme)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("failes to create a new ConfigMap with invalid owner", func() {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: namespace,
			},
		}

		invalidOwner := &metav1.ObjectMeta{
			Name:      "invalidOwner",
			Namespace: namespace,
		}

		result, err := utils.ReconcileConfigMap(context.Background(), configMap, invalidOwner, k8sClient, testEnv.Scheme)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeNil())
	})
})
