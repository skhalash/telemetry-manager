package otelcollector

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	testutils "github.com/kyma-project/telemetry-manager/internal/utils/test"
)

func TestGateway_ApplyResources(t *testing.T) {
	image := "opentelemetry/collector:latest"
	namespace := "kyma-system"
	priorityClassName := "normal"

	tests := []struct {
		name           string
		sut            *GatewayApplierDeleter
		goldenFilePath string
	}{
		{
			name:           "metric gateway",
			sut:            NewMetricGatewayApplierDeleter(image, namespace, priorityClassName),
			goldenFilePath: "testdata/metric-gateway.yaml",
		},
		{
			name:           "trace gateway",
			sut:            NewTraceGatewayApplierDeleter(image, namespace, priorityClassName),
			goldenFilePath: "testdata/trace-gateway.yaml",
		},
		{
			name:           "log gateway",
			sut:            NewLogGatewayApplierDeleter(image, namespace, priorityClassName),
			goldenFilePath: "testdata/log-gateway.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objects []client.Object

			client := fake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(_ context.Context, c client.WithWatch, obj client.Object, _ ...client.CreateOption) error {
					objects = append(objects, obj)
					// Nothing has to be created, just add created object to the list
					return nil
				},
			}).Build()

			err := tt.sut.ApplyResources(context.Background(), client, GatewayApplyOptions{
				AllowedPorts:        []int32{5555, 6666},
				CollectorConfigYAML: "dummy",
				CollectorEnvVars: map[string][]byte{
					"DUMMY_ENV_VAR": []byte("foo"),
				},
				Replicas: 2,
			})
			require.NoError(t, err)

			bytes, err := testutils.MarshalYAML(objects)
			require.NoError(t, err)

			goldenFileBytes, err := os.ReadFile(tt.goldenFilePath)
			require.NoError(t, err)

			require.Equal(t, string(goldenFileBytes), string(bytes))
		})
	}
}

func TestGateway_DeleteResources(t *testing.T) {
	image := "opentelemetry/collector:latest"
	namespace := "kyma-system"
	priorityClassName := "normal"

	tests := []struct {
		name string
		sut  *GatewayApplierDeleter
	}{
		{
			name: "metric gateway",
			sut:  NewMetricGatewayApplierDeleter(image, namespace, priorityClassName),
		},
		{
			name: "trace gateway",
			sut:  NewTraceGatewayApplierDeleter(image, namespace, priorityClassName),
		},
		{
			name: "log gateway",
			sut:  NewLogGatewayApplierDeleter(image, namespace, priorityClassName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var created []client.Object

			fakeClient := fake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(ctx context.Context, c client.WithWatch, obj client.Object, _ ...client.CreateOption) error {
					created = append(created, obj)
					return c.Create(ctx, obj)
				},
			}).Build()

			err := tt.sut.ApplyResources(context.Background(), fakeClient, GatewayApplyOptions{
				AllowedPorts:        []int32{5555, 6666},
				CollectorConfigYAML: "dummy",
			})
			require.NoError(t, err)

			err = tt.sut.DeleteResources(context.Background(), fakeClient, false)
			require.NoError(t, err)

			for i := range created {
				// an update operation on a non-existent object should return a NotFound error
				err = fakeClient.Get(context.Background(), client.ObjectKeyFromObject(created[i]), created[i])
				require.True(t, apierrors.IsNotFound(err), "want not found, got %v: %#v", err, created[i])
			}
		})
	}
}
