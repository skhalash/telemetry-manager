# LogPipeline

The `logpipeline.telemetry.kyma-project.io` CustomResourceDefinition (CRD) is a detailed description of the kind of data and the format used to filter and ship application logs in Kyma. To get the current CRD and show the output in the YAML format, run this command:

```bash
kubectl get crd logpipeline.telemetry.kyma-project.io -o yaml
```

## Sample Custom Resource

The following LogPipeline object defines a pipeline integrating with the HTTP/JSON-based output. It uses basic authentication and excludes application logs emitted by `istio-proxy` containers.

```yaml
apiVersion: telemetry.kyma-project.io/v1alpha1
kind: LogPipeline
metadata:
  name: custom-fluentd
  generation: 2
spec:
  input:
    application:
      containers:
        exclude:
        - istio-proxy
      namespaces: {}
  output:
    http:
      dedot: true
      host:
        valueFrom:
          secretKeyRef:
            key: Fluentd-endpoint
            name: custom-fluentd
            namespace: default
      password:
        valueFrom:
          secretKeyRef:
            key: Fluentd-password
            name: custom-fluentd
            namespace: default
      tls: {}
      uri: /customindex/kyma
      user:
        valueFrom:
          secretKeyRef:
            key: Fluentd-username
            name: custom-fluentd
            namespace: default
status:
  conditions:
  - lastTransitionTime: "2024-02-28T22:48:24Z"
    message: Fluent Bit DaemonSet is ready
    observedGeneration: 2
    reason: AgentReady
    status: "True"
    type: AgentHealthy
  - lastTransitionTime: "2024-02-28T22:48:11Z"
    message: ""
    observedGeneration: 2
    reason: ConfigurationGenerated
    status: "True"
    type: ConfigurationGenerated
```

For further examples, see the [samples](https://github.com/kyma-project/telemetry-manager/tree/main/config/samples) directory.

## Custom Resource Parameters

For details, see the [LogPipeline specification file](https://github.com/kyma-project/telemetry-manager/blob/main/apis/telemetry/v1alpha1/logpipeline_types.go).

<!-- The table below was generated automatically -->
<!-- Some special tags (html comments) are at the end of lines due to markdown requirements. -->
<!-- The content between "TABLE-START" and "TABLE-END" will be replaced -->

<!-- TABLE-START -->
### LogPipeline.telemetry.kyma-project.io/v1alpha1

**Spec:**

| Parameter | Type | Description |
| ---- | ----------- | ---- |
| **files**  | \[\]object | Provides file content to be consumed by a LogPipeline configuration |
| **files.&#x200b;content**  | string |  |
| **files.&#x200b;name**  | string |  |
| **filters**  | \[\]object | Describes a filtering option on the logs of the pipeline. |
| **filters.&#x200b;custom**  | string | Custom filter definition in the Fluent Bit syntax. Note: If you use a `custom` filter, you put the LogPipeline in unsupported mode. |
| **input**  | object | Defines where to collect logs, including selector mechanisms. |
| **input.&#x200b;application**  | object | Configures in more detail from which containers application logs are enabled as input. |
| **input.&#x200b;application.&#x200b;containers**  | object | Describes whether application logs from specific containers are selected. The options are mutually exclusive. |
| **input.&#x200b;application.&#x200b;containers.&#x200b;exclude**  | \[\]string | Specifies to exclude only the container logs with the specified container names. |
| **input.&#x200b;application.&#x200b;containers.&#x200b;include**  | \[\]string | Specifies to include only the container logs with the specified container names. |
| **input.&#x200b;application.&#x200b;dropLabels**  | boolean | Defines whether to drop all Kubernetes labels. The default is `false`. |
| **input.&#x200b;application.&#x200b;enabled**  | boolean | If enabled, application logs are collected. The default is `true`. |
| **input.&#x200b;application.&#x200b;keepAnnotations**  | boolean | Defines whether to keep all Kubernetes annotations. The default is `false`. |
| **input.&#x200b;application.&#x200b;keepOriginalBody**  | boolean | If the `log` attribute contains a JSON payload and it is successfully parsed, the `log` attribute will be retained if `KeepOriginalBody` is set to `true`. Otherwise, the log attribute will be removed from the log record. The default is `true`. |
| **input.&#x200b;application.&#x200b;namespaces**  | object | Describes whether application logs from specific Namespaces are selected. The options are mutually exclusive. System Namespaces are excluded by default from the collection. |
| **input.&#x200b;application.&#x200b;namespaces.&#x200b;exclude**  | \[\]string | Exclude the container logs of the specified Namespace names. |
| **input.&#x200b;application.&#x200b;namespaces.&#x200b;include**  | \[\]string | Include only the container logs of the specified Namespace names. |
| **input.&#x200b;application.&#x200b;namespaces.&#x200b;system**  | boolean | Set to `true` if collecting from all Namespaces must also include the system Namespaces like kube-system, istio-system, and kyma-system. |
| **input.&#x200b;otlp**  | object | Configures an endpoint to receive logs from a OTLP source. |
| **input.&#x200b;otlp.&#x200b;disabled**  | boolean | If set to `true`, no push-based OTLP signals are collected. The default is `false`. |
| **input.&#x200b;otlp.&#x200b;namespaces**  | object | Describes whether push-based OTLP signals from specific namespaces are selected. System namespaces are enabled by default. |
| **input.&#x200b;otlp.&#x200b;namespaces.&#x200b;exclude**  | \[\]string | Exclude signals from the specified Namespace names only. |
| **input.&#x200b;otlp.&#x200b;namespaces.&#x200b;include**  | \[\]string | Include signals from the specified Namespace names only. |
| **output**  | object | [Fluent Bit output](https://docs.fluentbit.io/manual/pipeline/outputs) where you want to push the logs. Only one output can be specified. |
| **output.&#x200b;custom**  | string | Defines a custom output in the Fluent Bit syntax. Note: If you use a `custom` output, you put the LogPipeline in unsupported mode. |
| **output.&#x200b;http**  | object | Configures an HTTP-based output compatible with the Fluent Bit HTTP output plugin. |
| **output.&#x200b;http.&#x200b;compress**  | string | Defines the compression algorithm to use. |
| **output.&#x200b;http.&#x200b;dedot**  | boolean | Enables de-dotting of Kubernetes labels and annotations for compatibility with ElasticSearch based backends. Dots (.) will be replaced by underscores (_). Default is `false`. |
| **output.&#x200b;http.&#x200b;format**  | string | Data format to be used in the HTTP request body. Default is `json`. |
| **output.&#x200b;http.&#x200b;host**  | object | Defines the host of the HTTP receiver. |
| **output.&#x200b;http.&#x200b;host.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;host.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;host.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;host.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;host.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;host.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;http.&#x200b;password**  | object | Defines the basic auth password. |
| **output.&#x200b;http.&#x200b;password.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;password.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;http.&#x200b;port**  | string | Defines the port of the HTTP receiver. Default is 443. |
| **output.&#x200b;http.&#x200b;tls**  | object | Configures TLS for the HTTP target server. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca**  | object | Defines an optional CA certificate for server certificate verification when using TLS. The certificate must be provided in PEM format. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert**  | object | Defines a client certificate to use when using TLS. The certificate must be provided in PEM format. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;disabled**  | boolean | Indicates if TLS is disabled or enabled. Default is `false`. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key**  | object | Defines the client key to use when using TLS. The key must be provided in PEM format. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;http.&#x200b;tls.&#x200b;skipCertificateValidation**  | boolean | If `true`, the validation of certificates is skipped. Default is `false`. |
| **output.&#x200b;http.&#x200b;uri**  | string | Defines the URI of the HTTP receiver. Default is "/". |
| **output.&#x200b;http.&#x200b;user**  | object | Defines the basic auth user. |
| **output.&#x200b;http.&#x200b;user.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;http.&#x200b;user.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;http.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;http.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;http.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;http.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp**  | object | Defines an output using the OpenTelemetry protocol. |
| **output.&#x200b;otlp.&#x200b;authentication**  | object | Defines authentication options for the OTLP output |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic**  | object | Activates `Basic` authentication for the destination providing relevant Secrets. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password** (required) | object | Contains the basic auth password or a Secret reference. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;password.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user** (required) | object | Contains the basic auth username or a Secret reference. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;authentication.&#x200b;basic.&#x200b;user.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;endpoint** (required) | object | Defines the host and port (<host>:<port>) of an OTLP endpoint. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;endpoint.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;headers**  | \[\]object | Defines custom headers to be added to outgoing HTTP or GRPC requests. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;name** (required) | string | Defines the header name. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;prefix**  | string | Defines an optional header value prefix. The prefix is separated from the value by a space character. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;headers.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;path**  | string | Defines OTLP export URL path (only for the HTTP protocol). This value overrides auto-appended paths `/v1/metrics` and `/v1/traces` |
| **output.&#x200b;otlp.&#x200b;protocol**  | string | Defines the OTLP protocol (http or grpc). Default is grpc. |
| **output.&#x200b;otlp.&#x200b;tls**  | object | Defines TLS options for the OTLP output. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca**  | object | Defines an optional CA certificate for server certificate verification when using TLS. The certificate must be provided in PEM format. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;ca.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert**  | object | Defines a client certificate to use when using TLS. The certificate must be provided in PEM format. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;cert.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;insecure**  | boolean | Defines whether to send requests using plaintext instead of TLS. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;insecureSkipVerify**  | boolean | Defines whether to skip server certificate verification when using TLS. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key**  | object | Defines the client key to use when using TLS. The key must be provided in PEM format. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;value**  | string | The value as plain text. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;valueFrom**  | object | The value as a reference to a resource. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **output.&#x200b;otlp.&#x200b;tls.&#x200b;key.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |
| **variables**  | \[\]object | A list of mappings from Kubernetes Secret keys to environment variables. Mapped keys are mounted as environment variables, so that they are available as [Variables](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/classic-mode/variables) in the sections. |
| **variables.&#x200b;name**  | string | Name of the variable to map. |
| **variables.&#x200b;valueFrom**  | object |  |
| **variables.&#x200b;valueFrom.&#x200b;secretKeyRef**  | object | Refers to the value of a specific key in a Secret. You must provide `name` and `namespace` of the Secret, as well as the name of the `key`. |
| **variables.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;key** (required) | string | The name of the attribute of the Secret holding the referenced value. |
| **variables.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;name** (required) | string | The name of the Secret containing the referenced value. |
| **variables.&#x200b;valueFrom.&#x200b;secretKeyRef.&#x200b;namespace** (required) | string | The name of the namespace containing the Secret with the referenced value. |

**Status:**

| Parameter | Type | Description |
| ---- | ----------- | ---- |
| **conditions**  | \[\]object | An array of conditions describing the status of the pipeline. |
| **conditions.&#x200b;lastTransitionTime** (required) | string | lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable. |
| **conditions.&#x200b;message** (required) | string | message is a human readable message indicating details about the transition. This may be an empty string. |
| **conditions.&#x200b;observedGeneration**  | integer | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. |
| **conditions.&#x200b;reason** (required) | string | reason contains a programmatic identifier indicating the reason for the condition's last transition. Producers of specific condition types may define expected values and meanings for this field, and whether the values are considered a guaranteed API. The value should be a CamelCase string. This field may not be empty. |
| **conditions.&#x200b;status** (required) | string | status of the condition, one of True, False, Unknown. |
| **conditions.&#x200b;type** (required) | string | type of condition in CamelCase or in foo.example.com/CamelCase. |
| **unsupportedMode**  | boolean | Is active when the LogPipeline uses a `custom` output or filter; see [unsupported mode](https://github.com/kyma-project/telemetry-manager/blob/main/docs/user/02-logs.md#unsupported-mode). |

<!-- TABLE-END -->

### LogPipeline Status

The status of the LogPipeline is determined by the condition types `AgentHealthy`, `ConfigurationGenerated`, and `TelemetryFlowHealthy`:

| Condition Type         | Condition Status | Condition Reason             | Condition Message                                                                                                                                                                                                                  |
|------------------------|------------------|------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| AgentHealthy           | True             | AgentReady                   | Log agent DaemonSet is ready                                                                                                                                                                                                       |
| AgentHealthy           | True             | RolloutInProgress            | Pods are being started/updated                                                                                                                                                                                                     |
| AgentHealthy           | False            | AgentNotReady                | No Pods deployed                                                                                                                                                                                                                   |
| AgentHealthy           | False            | AgentNotReady                | DaemonSet is not yet created                                                                                                                                                                                                       |
| AgentHealthy           | False            | AgentNotReady                | Failed to get DaemonSet                                                                                                                                                                                                            |
| AgentHealthy           | False            | AgentNotReady                | Pod is in the pending state because container: `container name` is not running due to: `reason`. Please check the container: `container name` logs.                                                                                |
| AgentHealthy           | False            | AgentNotReady                | Pod is in the failed state due to: `reason`                                                                                                                                                                                        |
| ConfigurationGenerated | True             | AgentConfigured              | LogPipeline specification is successfully applied to the configuration of Fluent Bit agent                                                                                                                                         |
| ConfigurationGenerated | True             | TLSCertificateAboutToExpire  | TLS (CA) certificate is about to expire, configured certificate is valid until YYYY-MM-DD                                                                                                                                          |
| ConfigurationGenerated | False            | EndpointInvalid              | HTTP output host invalid: `reason`                                                                                                                                                                                                 |
| ConfigurationGenerated | False            | ReferencedSecretMissing      | One or more referenced Secrets are missing: Secret 'my-secret' of Namespace 'my-namespace'                                                                                                                                         |
| ConfigurationGenerated | False            | ReferencedSecretMissing      | One or more keys in a referenced Secret are missing: Key 'my-key' in Secret 'my-secret' of Namespace 'my-namespace'"                                                                                                               |
| ConfigurationGenerated | False            | ReferencedSecretMissing      | Secret reference is missing field/s: (field1, field2, ...)                                                                                                                                                                         |
| ConfigurationGenerated | False            | TLSCertificateExpired        | TLS (CA) certificate expired on YYYY-MM-DD                                                                                                                                                                                         |
| ConfigurationGenerated | False            | TLSConfigurationInvalid      | TLS configuration invalid                                                                                                                                                                                                          |
| ConfigurationGenerated | False            | ValidationFailed             | Pipeline validation failed due to an error from the Kubernetes API server                                                                                                                                                          |
| TelemetryFlowHealthy   | True             | FlowHealthy                  | No problems detected in the telemetry flow                                                                                                                                                                                         |
| TelemetryFlowHealthy   | False            | AgentAllTelemetryDataDropped | Backend is not reachable or rejecting logs. All logs are dropped. See troubleshooting: [No Logs Arrive at the Backend](https://kyma-project.io/#/telemetry-manager/user/02-logs?id=no-logs-arrive-at-the-backend)                  |
| TelemetryFlowHealthy   | False            | AgentBufferFillingUp         | Buffer nearing capacity. Incoming log rate exceeds export rate. See troubleshooting: [Agent Buffer Filling Up](https://kyma-project.io/#/telemetry-manager/user/02-logs?id=agent-buffer-filling-up)                                |
| TelemetryFlowHealthy   | False            | AgentNoLogsDelivered         | Backend is not reachable or rejecting logs. Logs are buffered and not yet dropped. See troubleshooting: [No Logs Arrive at the Backend](https://kyma-project.io/#/telemetry-manager/user/02-logs?id=no-logs-arrive-at-the-backend) |
| TelemetryFlowHealthy   | False            | AgentSomeDataDropped         | Backend is reachable, but rejecting logs. Some logs are dropped. See troubleshooting: [Not All Logs Arrive at the Backend](https://kyma-project.io/#/telemetry-manager/user/02-logs?id=not-all-logs-arrive-at-the-backend)         |
| TelemetryFlowHealthy   | False            | ConfigurationNotGenerated    | No logs delivered to backend because LogPipeline specification is not applied to the configuration of Fluent Bit agent. Check the 'ConfigurationGenerated' condition for more details                                              |
| TelemetryFlowHealthy   | Unknown          | AgentProbingFailed               | Could not determine the health of the telemetry flow because the self monitor probing failed                                                                                                                                       |
