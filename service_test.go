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
	serviceName = "sample-service"
)

var _ = Describe("Reconcile Service", func() {
	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &corev1.Service{}, client.InNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when owner is valid", func() {
		It("creates a new Service", func() {
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceName,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 3000,
						},
					},
				},
			}

			result, err := utils.ReconcileService(context.Background(), service, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})

		It("updates an outdated Service", func() {
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceName,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 3000,
						},
					},
				},
			}

			result, err := utils.ReconcileService(context.Background(), service, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))

			updatedService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceName,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 3001,
						},
					},
				},
			}

			result, err = utils.ReconcileService(context.Background(), updatedService, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})
	})

	Context("when owner is invalid", func() {
		It("failes to create a new Service", func() {
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceName,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 3000,
						},
					},
				},
			}

			invalidOwner := &metav1.ObjectMeta{
				Name:      "invalidOwner",
				Namespace: namespace,
			}

			result, err := utils.ReconcileService(context.Background(), service, invalidOwner, k8sClient, testEnv.Scheme)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(utils.ReconcileResult{}))
		})
	})
})
