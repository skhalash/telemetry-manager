package otelcollector

import (
	"testing"

	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
)

func TestMakeTraceGatewayRBAC(t *testing.T) {
	namespace := "test-namespace"
	rbac := makeTraceGatewayRBAC(namespace)

	expectedName := TraceGatewayName
	t.Run("should have a cluster role", func(t *testing.T) {
		cr := rbac.clusterRole
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"replicasets"},
				Verbs:     []string{"get", "list", "watch"},
			},
		}

		require.NotNil(t, cr)
		require.Equal(t, expectedName, cr.Name)
		require.Equal(t, namespace, cr.Namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, cr.Labels)
		require.Equal(t, expectedRules, cr.Rules)
	})

	t.Run("should have a cluster role binding", func(t *testing.T) {
		crb := rbac.clusterRoleBinding
		checkClusterRoleBinding(t, crb, expectedName, namespace)
	})

	t.Run("should not have a role", func(t *testing.T) {
		r := rbac.role
		require.Nil(t, r)
	})

	t.Run("should not have a role binding", func(t *testing.T) {
		rb := rbac.roleBinding
		require.Nil(t, rb)
	})
}

func TestMakeMetricAgentRBAC(t *testing.T) {
	namespace := "test-namespace"
	rbac := makeMetricAgentRBAC(namespace)

	expectedName := "test-agent"
	t.Run("should have a cluster role", func(t *testing.T) {
		cr := rbac.clusterRole
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "nodes/stats", "nodes/proxy"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "nodes/metrics", "services", "endpoints", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				NonResourceURLs: []string{"/metrics", "/metrics/cadvisor"},
				Verbs:           []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events", "namespaces", "namespaces/status", "nodes", "nodes/spec", "pods", "pods/status", "replicationcontrollers", "replicationcontrollers/status", "resourcequotas", "services"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{"apps"},
				Resources: []string{"daemonsets", "deployments", "replicasets", "statefulsets"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{"extensions"},
				Resources: []string{"daemonsets", "deployments", "replicasets"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"get", "list", "watch"},
			},
		}

		require.NotNil(t, cr)
		require.Equal(t, cr.Name, expectedName)
		require.Equal(t, cr.Namespace, namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, cr.Labels)
		require.Equal(t, cr.Rules, expectedRules)
	})

	t.Run("should have a cluster role binding", func(t *testing.T) {
		crb := rbac.clusterRoleBinding
		checkClusterRoleBinding(t, crb, expectedName, namespace)
	})

	t.Run("should have a role", func(t *testing.T) {
		r := rbac.role
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{"coordination.k8s.io"},
				Resources: []string{"leases"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			},
		}

		require.NotNil(t, r)
		require.Equal(t, expectedName, r.Name)
		require.Equal(t, namespace, r.Namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, r.Labels)
		require.Equal(t, expectedRules, r.Rules)
	})

	t.Run("should have a role binding", func(t *testing.T) {
		rb := rbac.roleBinding
		require.NotNil(t, rb)

		checkRoleBinding(t, rb, expectedName, namespace)
	})
}

func TestMakeMetricGatewayRBAC(t *testing.T) {
	namespace := "test-namespace"
	rbac := makeMetricGatewayRBAC(namespace)

	expectedName := MetricGatewayName
	t.Run("should have a cluster role", func(t *testing.T) {
		cr := rbac.clusterRole
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"replicasets"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"operator.kyma-project.io"},
				Resources: []string{"telemetries"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"telemetry.kyma-project.io"},
				Resources: []string{"metricpipelines"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"telemetry.kyma-project.io"},
				Resources: []string{"tracepipelines"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"telemetry.kyma-project.io"},
				Resources: []string{"logpipelines"},
				Verbs:     []string{"get", "list", "watch"},
			}}

		require.NotNil(t, cr)
		require.Equal(t, expectedName, cr.Name)
		require.Equal(t, namespace, cr.Namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, cr.Labels)
		require.Equal(t, expectedRules, cr.Rules)
	})

	t.Run("should have a cluster role binding", func(t *testing.T) {
		crb := rbac.clusterRoleBinding
		checkClusterRoleBinding(t, crb, expectedName, namespace)
	})

	t.Run("should have a role", func(t *testing.T) {
		r := rbac.role
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{"coordination.k8s.io"},
				Resources: []string{"leases"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			},
		}

		require.NotNil(t, r)
		require.Equal(t, expectedName, r.Name)
		require.Equal(t, namespace, r.Namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, r.Labels)
		require.Equal(t, expectedRules, r.Rules)
	})

	t.Run("should have a role binding", func(t *testing.T) {
		rb := rbac.roleBinding
		require.NotNil(t, rb)

		checkRoleBinding(t, rb, expectedName, namespace)
	})
}

func TestMakeLogGatewayRBAC(t *testing.T) {
	namespace := "test-namespace"
	rbac := makeLogGatewayRBAC(namespace)

	expectedName := LogGatewayName
	t.Run("should have a cluster role", func(t *testing.T) {
		cr := rbac.clusterRole
		expectedRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"replicasets"},
				Verbs:     []string{"get", "list", "watch"},
			},
		}

		require.NotNil(t, cr)
		require.Equal(t, expectedName, cr.Name)
		require.Equal(t, namespace, cr.Namespace)
		require.Equal(t, map[string]string{
			"app.kubernetes.io/name": expectedName,
		}, cr.Labels)
		require.Equal(t, expectedRules, cr.Rules)
	})

	t.Run("should have a cluster role binding", func(t *testing.T) {
		crb := rbac.clusterRoleBinding
		checkClusterRoleBinding(t, crb, expectedName, namespace)
	})

	t.Run("should not have a role", func(t *testing.T) {
		r := rbac.role
		require.Nil(t, r)
	})

	t.Run("should not have a role binding", func(t *testing.T) {
		rb := rbac.roleBinding
		require.Nil(t, rb)
	})
}

func checkClusterRoleBinding(t *testing.T, crb *rbacv1.ClusterRoleBinding, name, namespace string) {
	require.NotNil(t, crb)
	require.Equal(t, name, crb.Name)
	require.Equal(t, namespace, crb.Namespace)
	require.Equal(t, map[string]string{
		"app.kubernetes.io/name": name,
	}, crb.Labels)

	subject := crb.Subjects[0]
	require.Equal(t, "ServiceAccount", subject.Kind)
	require.Equal(t, name, subject.Name)
	require.Equal(t, namespace, subject.Namespace)

	require.Equal(t, "rbac.authorization.k8s.io", crb.RoleRef.APIGroup)
	require.Equal(t, "ClusterRole", crb.RoleRef.Kind)
	require.Equal(t, name, crb.RoleRef.Name)
}

func checkRoleBinding(t *testing.T, rb *rbacv1.RoleBinding, name, namespace string) {
	require.Equal(t, name, rb.Name)
	require.Equal(t, namespace, rb.Namespace)
	require.Equal(t, map[string]string{
		"app.kubernetes.io/name": name,
	}, rb.Labels)

	subject := rb.Subjects[0]
	require.Equal(t, "ServiceAccount", subject.Kind)
	require.Equal(t, name, subject.Name)
	require.Equal(t, namespace, subject.Namespace)

	require.Equal(t, "rbac.authorization.k8s.io", rb.RoleRef.APIGroup)
	require.Equal(t, "Role", rb.RoleRef.Kind)
	require.Equal(t, name, rb.RoleRef.Name)
}
