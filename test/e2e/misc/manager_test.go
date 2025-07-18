//go:build e2e

package misc

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kitkyma "github.com/kyma-project/telemetry-manager/test/testkit/kyma"
	"github.com/kyma-project/telemetry-manager/test/testkit/periodic"
	"github.com/kyma-project/telemetry-manager/test/testkit/suite"
)

var _ = Describe(suite.ID(), func() {
	Context("After deploying manifest", func() {
		It("Should have kyma-system namespace", Label(suite.LabelTelemetry), func() {
			var namespace corev1.Namespace
			key := types.NamespacedName{
				Name: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &namespace)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have a running manager deployment", Label(suite.LabelTelemetry), func() {
			var deployment appsv1.Deployment
			key := types.NamespacedName{
				Name:      "telemetry-manager",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &deployment)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				listOptions := client.ListOptions{
					LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels),
					Namespace:     kitkyma.SystemNamespaceName,
				}
				var pods corev1.PodList
				err := suite.K8sClient.List(suite.Ctx, &pods, &listOptions)
				Expect(err).NotTo(HaveOccurred())
				for _, pod := range pods.Items {
					for _, containerStatus := range pod.Status.ContainerStatuses {
						if containerStatus.State.Running == nil {
							return false
						}
					}
				}

				return true
			}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(BeTrue())
		})

		It("Should have a webhook service", Label(suite.LabelTelemetry), func() {
			var service corev1.Service
			err := suite.K8sClient.Get(suite.Ctx, kitkyma.TelemetryManagerWebhookServiceName, &service)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() []string {
				var endpointsList discoveryv1.EndpointSliceList
				err := suite.K8sClient.List(suite.Ctx, &endpointsList, client.InNamespace(kitkyma.SystemNamespaceName))
				Expect(err).NotTo(HaveOccurred())

				var webhookEndpoints *discoveryv1.EndpointSlice
				for _, endpoints := range endpointsList.Items {
					// EndpointSlice names are prefixed with the service name
					if strings.HasPrefix(endpoints.Name, kitkyma.TelemetryManagerWebhookServiceName.Name) {
						webhookEndpoints = &endpoints
						break
					}
				}
				Expect(webhookEndpoints).NotTo(BeNil())

				var addresses []string
				for _, endpoint := range webhookEndpoints.Endpoints {
					addresses = append(addresses, endpoint.Addresses...)
				}
				return addresses
			}, periodic.EventuallyTimeout, periodic.DefaultInterval).ShouldNot(BeEmpty(), "Webhook service endpoints should have IP addresses assigned")
		})

		It("Should have a metrics service", Label(suite.LabelTelemetry), func() {
			var service corev1.Service
			err := suite.K8sClient.Get(suite.Ctx, kitkyma.TelemetryManagerMetricsServiceName, &service)
			Expect(err).NotTo(HaveOccurred())

			Expect(service.Annotations).Should(HaveKeyWithValue("prometheus.io/scrape", "true"))
			Expect(service.Annotations).Should(HaveKeyWithValue("prometheus.io/port", "8080"))

			Eventually(func() []string {
				var endpointsList discoveryv1.EndpointSliceList
				err := suite.K8sClient.List(suite.Ctx, &endpointsList, client.InNamespace(kitkyma.SystemNamespaceName))
				Expect(err).NotTo(HaveOccurred())

				var metricsEndpoints *discoveryv1.EndpointSlice
				for _, endpoints := range endpointsList.Items {
					// EndpointSlice names are prefixed with the service name
					if strings.HasPrefix(endpoints.Name, kitkyma.TelemetryManagerMetricsServiceName.Name) {
						metricsEndpoints = &endpoints
						break
					}
				}
				Expect(metricsEndpoints).NotTo(BeNil())

				var addresses []string
				for _, endpoint := range metricsEndpoints.Endpoints {
					addresses = append(addresses, endpoint.Addresses...)
				}
				return addresses
			}, periodic.EventuallyTimeout, periodic.DefaultInterval).ShouldNot(BeEmpty(), "Metrics service endpoints should have IP addresses assigned")
		})

		It("Should have LogPipelines CRD", Label(suite.LabelFluentBit), func() {
			var crd apiextensionsv1.CustomResourceDefinition
			key := types.NamespacedName{
				Name: "logpipelines.telemetry.kyma-project.io",
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &crd)
			Expect(err).NotTo(HaveOccurred())
			Expect(crd.Spec.Scope).To(Equal(apiextensionsv1.ClusterScoped))
		})

		It("Should have LogParsers CRD", Label(suite.LabelFluentBit), func() {
			var crd apiextensionsv1.CustomResourceDefinition
			key := types.NamespacedName{
				Name: "logparsers.telemetry.kyma-project.io",
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &crd)
			Expect(err).NotTo(HaveOccurred())
			Expect(crd.Spec.Scope).To(Equal(apiextensionsv1.ClusterScoped))
		})

		It("Should have TracePipelines CRD", Label(suite.LabelTraces), func() {
			var crd apiextensionsv1.CustomResourceDefinition
			key := types.NamespacedName{
				Name: "tracepipelines.telemetry.kyma-project.io",
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &crd)
			Expect(err).NotTo(HaveOccurred())
			Expect(crd.Spec.Scope).To(Equal(apiextensionsv1.ClusterScoped))
		})

		It("Should have MetricPipelines CRD", Label(suite.LabelMetrics), func() {
			var crd apiextensionsv1.CustomResourceDefinition
			key := types.NamespacedName{
				Name: "metricpipelines.telemetry.kyma-project.io",
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &crd)
			Expect(err).NotTo(HaveOccurred())
			Expect(crd.Spec.Scope).To(Equal(apiextensionsv1.ClusterScoped))
		})

		It("Should have Telemetry CRD", Label(suite.LabelTelemetry), func() {
			var crd apiextensionsv1.CustomResourceDefinition
			key := types.NamespacedName{
				Name: "telemetries.operator.kyma-project.io",
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &crd)
			Expect(err).NotTo(HaveOccurred())
			Expect(crd.Spec.Scope).To(Equal(apiextensionsv1.NamespaceScoped))
		})

		It("Should have a Busola extension for MetricPipelines CRD", Label(suite.LabelMetrics), func() {
			var cm corev1.ConfigMap
			key := types.NamespacedName{
				Name:      "telemetry-metricpipelines",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &cm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have a Busola extension for LogPipelines CRD", Label(suite.LabelFluentBit), func() {
			var cm corev1.ConfigMap
			key := types.NamespacedName{
				Name:      "telemetry-logpipelines",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &cm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have a Busola extension for TracePipelines CRD", Label(suite.LabelTraces), func() {
			var cm corev1.ConfigMap
			key := types.NamespacedName{
				Name:      "telemetry-tracepipelines",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &cm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have a Busola extension for Telemetry CRD", Label(suite.LabelTelemetry), func() {
			var cm corev1.ConfigMap
			key := types.NamespacedName{
				Name:      "telemetry-module",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &cm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have a NetworkPolicy", Label(suite.LabelTelemetry), func() {
			var networkPolicy networkingv1.NetworkPolicy
			key := types.NamespacedName{
				Name:      "telemetry-manager",
				Namespace: kitkyma.SystemNamespaceName,
			}
			err := suite.K8sClient.Get(suite.Ctx, key, &networkPolicy)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have priority class resource created", Label(suite.LabelTelemetry), func() {
			priorityClassNames := []string{"telemetry-priority-class", "telemetry-priority-class-high"}
			var priorityClass schedulingv1.PriorityClass
			for _, prioClass := range priorityClassNames {
				key := types.NamespacedName{
					Name:      prioClass,
					Namespace: kitkyma.SystemNamespaceName,
				}
				err := suite.K8sClient.Get(suite.Ctx, key, &priorityClass)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
