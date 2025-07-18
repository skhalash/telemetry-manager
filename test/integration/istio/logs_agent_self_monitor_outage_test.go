//go:build istio

package istio

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorv1alpha1 "github.com/kyma-project/telemetry-manager/apis/operator/v1alpha1"
	"github.com/kyma-project/telemetry-manager/internal/conditions"
	testutils "github.com/kyma-project/telemetry-manager/internal/utils/test"
	"github.com/kyma-project/telemetry-manager/test/testkit/assert"
	kitk8s "github.com/kyma-project/telemetry-manager/test/testkit/k8s"
	kitkyma "github.com/kyma-project/telemetry-manager/test/testkit/kyma"
	. "github.com/kyma-project/telemetry-manager/test/testkit/matchers/prometheus"
	kitbackend "github.com/kyma-project/telemetry-manager/test/testkit/mocks/backend"
	"github.com/kyma-project/telemetry-manager/test/testkit/mocks/floggen"
	"github.com/kyma-project/telemetry-manager/test/testkit/suite"
)

var _ = Describe(suite.ID(), Label(suite.LabelSelfMonitoringLogsAgentOutage), Ordered, func() {
	var (
		mockNs       = "istio-permissive-mtls"
		pipelineName = suite.ID()
		backend      *kitbackend.Backend
	)

	makeResources := func() []client.Object {

		backend = kitbackend.New(mockNs, kitbackend.SignalTypeLogsOTel, kitbackend.WithReplicas(0))

		logProducer := floggen.NewDeployment(mockNs).WithReplicas(3)

		logPipeline := testutils.NewLogPipelineBuilder().
			WithName(pipelineName).
			WithInput(testutils.BuildLogPipelineApplicationInput(testutils.ExtIncludeNamespaces(mockNs))).
			WithOTLPOutput(testutils.OTLPEndpoint(backend.Endpoint())).
			Build()

		objs := []client.Object{
			logProducer.K8sObject(),
			&logPipeline,
		}
		objs = append(objs, backend.K8sObjects()...)

		return objs
	}

	Context("When a logpipeline exists", Ordered, func() {
		BeforeAll(func() {
			k8sObjects := makeResources()
			DeferCleanup(func() {
				Expect(kitk8s.DeleteObjects(k8sObjects...)).Should(Succeed())
			})
			Expect(kitk8s.CreateObjects(GinkgoT(), k8sObjects...)).Should(Succeed())
		})

		It("Should have a running logpipeline", func() {
			assert.OTelLogPipelineHealthy(GinkgoT(), pipelineName)
		})

		It("Should have a running log agent daemonset", func() {
			assert.DaemonSetReady(GinkgoT(), kitkyma.LogAgentName)
		})

		It("Should have a running self-monitor", func() {
			assert.DeploymentReady(GinkgoT(), kitkyma.SelfMonitorName)
		})

		It("Should have a log producer running", func() {
			assert.DeploymentReady(GinkgoT(), types.NamespacedName{Namespace: mockNs, Name: floggen.DefaultName})
		})

		It("Should wait for the log flow to gradually become unhealthy", func() {
			assert.LogPipelineConditionReasonsTransition(GinkgoT(), pipelineName, conditions.TypeFlowHealthy, []assert.ReasonStatus{
				{Reason: conditions.ReasonSelfMonFlowHealthy, Status: metav1.ConditionTrue},
				{Reason: conditions.ReasonSelfMonAgentBufferFillingUp, Status: metav1.ConditionFalse},
				{Reason: conditions.ReasonSelfMonAgentAllDataDropped, Status: metav1.ConditionFalse},
			})

			assert.TelemetryHasState(GinkgoT(), operatorv1alpha1.StateWarning)
			assert.TelemetryHasCondition(GinkgoT(), suite.K8sClient, metav1.Condition{
				Type:   conditions.TypeLogComponentsHealthy,
				Status: metav1.ConditionFalse,
				Reason: conditions.ReasonSelfMonAgentAllDataDropped,
			})
		})

		Context("Metric instrumentation", Ordered, func() {
			It("Ensures that controller_runtime_webhook_requests_total is increased", func() {
				// Pushing metrics to the metric gateway triggers an alert.
				// It makes the self-monitor call the webhook, which in turn increases the counter.
				assert.EmitsManagerMetrics(GinkgoT(),
					HaveName(Equal("controller_runtime_webhook_requests_total")),
					SatisfyAll(
						HaveLabels(HaveKeyWithValue("webhook", "/api/v2/alerts")),
						HaveMetricValue(BeNumerically(">", 0)),
					))
			})

			It("Ensures that telemetry_self_monitor_prober_requests_total is emitted", func() {
				assert.EmitsManagerMetrics(GinkgoT(),
					HaveName(Equal("telemetry_self_monitor_prober_requests_total")),
					HaveMetricValue(BeNumerically(">", 0)),
				)
			})
		})
	})
})
