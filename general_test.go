package utils_test

import (
	"context"

	utils "github.com/cedrichopf/operator-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test DeleteResources", func() {
	Context("when a valid object will be deleted", func() {
		It("deletes the object from Kubernetes", func() {
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-configmap",
					Namespace: namespace,
				},
				Data: map[string]string{
					"sample": "data",
				},
			}

			err := k8sClient.Create(context.Background(), configMap)
			Expect(err).ToNot(HaveOccurred())

			existingConfigMap := &corev1.ConfigMap{}
			err = k8sClient.Get(context.Background(), types.NamespacedName{Name: "example-configmap", Namespace: namespace}, existingConfigMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(existingConfigMap).ToNot(Equal(&corev1.ConfigMap{}))

			err = utils.DeleteResources(context.Background(), k8sClient, []client.Object{existingConfigMap})
			Expect(err).ToNot(HaveOccurred())

			deletedConfigMap := &corev1.ConfigMap{}
			err = k8sClient.Get(context.Background(), types.NamespacedName{Name: "example-configmap", Namespace: namespace}, deletedConfigMap)
			Expect(apierrors.IsNotFound(err)).To(Equal(true))
		})
	})
})

var _ = Describe("Test RemoveFinalizer", func() {
	Context("when finalizer exists", func() {
		It("removes the finalizer", func() {
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "example-configmap",
					Namespace:  namespace,
					Finalizers: []string{"test-finalizer"},
				},
				Data: map[string]string{
					"sample": "data",
				},
			}

			removed := utils.RemoveFinalizer(configMap, "test-finalizer")
			Expect(removed).To(Equal(true))
		})
	})

	Context("when finalizer does not exist", func() {
		It("does not remove the finalizer", func() {
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "example-configmap",
					Namespace:  namespace,
					Finalizers: []string{"test-finalizer"},
				},
				Data: map[string]string{
					"sample": "data",
				},
			}

			removed := utils.RemoveFinalizer(configMap, "non-existing-finalizer")
			Expect(removed).To(Equal(false))
		})
	})
})
