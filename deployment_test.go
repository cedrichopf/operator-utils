package utils_test

import (
	"context"

	utils "github.com/cedrichopf/operator-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deploymentName = "sample-deployment"
)

var _ = Describe("Deployment", func() {
	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &appsv1.Deployment{}, client.InNamespace(namespace))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when owner is valid", func() {
		It("creates a new Deployment", func() {
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"test": "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"test": "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx",
								},
							},
						},
					},
				},
			}

			result, err := utils.ReconcileDeployment(context.Background(), deployment, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeNil())
		})

		It("updates an outdated Deployment", func() {
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"test": "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"test": "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx",
								},
							},
						},
					},
				},
			}

			result, err := utils.ReconcileDeployment(context.Background(), deployment, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeNil())

			updatedDeployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"test": "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"test": "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx-2",
									Image: "nginx",
								},
							},
						},
					},
				},
			}

			result, err = utils.ReconcileDeployment(context.Background(), updatedDeployment, owner, k8sClient, testEnv.Scheme)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Context("when owner is invalid", func() {
		It("failes to create a new Deployment", func() {
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: namespace,
				},
			}

			invalidOwner := &metav1.ObjectMeta{
				Name:      "invalidOwner",
				Namespace: namespace,
			}

			result, err := utils.ReconcileDeployment(context.Background(), deployment, invalidOwner, k8sClient, testEnv.Scheme)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})
})
