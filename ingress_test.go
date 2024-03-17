package utils_test

import (
	"context"

	utils "github.com/cedrichopf/operator-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingressName = "sample-ingress"
)

var _ = Describe("Reconcile Ingress", func() {
	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &networkingv1.Ingress{}, client.InNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when owner is valid", func() {
		It("creates a new Ingress", func() {
			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ingressName,
					Namespace: namespace,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "test.example.com",
						},
					},
				},
			}

			result, err := utils.ReconcileIngress(context.Background(), ingress, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})

		It("updates an outdated Ingress", func() {
			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ingressName,
					Namespace: namespace,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "test.example.com",
						},
					},
				},
			}

			result, err := utils.ReconcileIngress(context.Background(), ingress, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))

			updatedIngress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ingressName,
					Namespace: namespace,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "test-v2.example.com",
						},
					},
				},
			}

			result, err = utils.ReconcileIngress(context.Background(), updatedIngress, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Updated).To(Equal(true))
		})
	})

	Context("when owner is invalid", func() {
		It("failes to create a new Ingress", func() {
			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ingressName,
					Namespace: namespace,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "test.example.com",
						},
					},
				},
			}

			invalidOwner := &metav1.ObjectMeta{
				Name:      "invalidOwner",
				Namespace: namespace,
			}

			result, err := utils.ReconcileIngress(context.Background(), ingress, invalidOwner, k8sClient, testEnv.Scheme)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(utils.ReconcileResult{}))
		})
	})
})
