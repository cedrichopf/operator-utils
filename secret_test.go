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
	secretName = "sample-secret"
)

var _ = Describe("Reconcile Secret", func() {
	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &corev1.Secret{}, client.InNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when owner is valid", func() {
		It("creates a new Secret", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
			}

			result, err := utils.ReconcileSecret(context.Background(), secret, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})

		It("updates an outdated Secret", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
				StringData: map[string]string{
					"sample": "data",
				},
			}

			result, err := utils.ReconcileSecret(context.Background(), secret, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))

			updatedSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
				StringData: map[string]string{
					"sample": "newData",
				},
			}

			result, err = utils.ReconcileSecret(context.Background(), updatedSecret, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})
	})

	Context("when owner is invalid", func() {
		It("failes to create a new Secret", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
			}

			invalidOwner := &metav1.ObjectMeta{
				Name:      "invalidOwner",
				Namespace: namespace,
			}

			result, err := utils.ReconcileSecret(context.Background(), secret, invalidOwner, k8sClient, testEnv.Scheme)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(utils.ReconcileResult{}))
		})
	})
})
