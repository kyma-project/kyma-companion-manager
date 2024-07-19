# Custom Resource Definitions (CRDs)

## Companion Manager CR

The Companion Manager CR (Custom Resource) is the companion backend in Kyma.

### Specification

The Companion Manager CR has the following fields:

- `name`: The name of the companion application.
- `namespace`: The namespace in which the companion application is deployed.
- `deplyomentName`: The name of the deployment for the companion application.
- `replicas`: The number of replicas for the companion application.
- `image`: The Docker image of the companion application.
- `imagePullSecrets`: The name of the secret to be used for pulling the Docker image.
- `imagePullPolicy`: The policy for pulling the Docker image.
- `ports`: The ports to be exposed by the companion application.
- `resources`: The resource limits and requests for the companion application. (CPU and memory)
- `env`: Environment variables to be passed to the companion application. Source are configMap and secret.
- `labels`: Labels to be applied to the companion application.
- `annotations`: Annotations to be applied to the companion application.
- `serviceAccount`: The service account to be used by the companion application.
- `serviceAccountName`: The name of the service account to be used by the companion application.
- `restartPolicy`: The restart policy for the companion application.

From configuration perspective, not all fields are mandatory. The only mandatory fields depend on the following conditions:

- Flexibility
- Easy to deploy
- Easy to maintain
- Automation support

From this reason we should define the best Custom Resource Definition (CRD) for the Companion Manager CR.

### Options

#### 1. All fields are hardcoded in the CRD.

The user cannot change any field.

[config/crd/bases/all-hardcoded-operator.kyma-project.io_companions.yaml](config/crd/bases/all-hardcoded-operator.kyma-project.io_companions.yaml)

```yaml
properties:
```

#### 2. Namespaces are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace.

[config/crd/bases/namespaces-operator.kyma-project.io_companions.yaml](config/crd/bases/namespaces-operator.kyma-project.io_companions.yaml)

```yaml
properties:
  deploymentNamespace:
    type: string
  namespaces:
    items:
      type: string
    type: array
required:
  - deploymentNamespace
  - namespaces
type: object
```

#### 3. Namespaces, ConfigMaps, and Secrets are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, and Secrets.

[config/crd/bases/configmaps-secrets-operator.kyma-project.io_companions.yaml](config/crd/bases/configmaps-secrets-operator.kyma-project.io_companions.yaml)

```yaml
properties:
  configMapNames:
    items:
      type: string
    type: array
  deploymentNamespace:
    type: string
  namespaces:
    items:
      type: string
    type: array
  secretNames:
    items:
      type: string
    type: array
required:
  - configMapNames
  - deploymentNamespace
  - namespaces
  - secretNames
type: object
```

#### 4. Namespaces, ConfigMaps, Secrets and ContainerPort are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, and ContainerPort.

[config/crd/bases/containerport-operator.kyma-project.io_companions.yaml](config/crd/bases/containerport-operator.kyma-project.io_companions.yaml)

```yaml
properties:
  configMapNames:
    items:
      type: string
    type: array
  containerPort:
    format: int32
    type: integer
  deploymentNamespace:
    type: string
  namespaces:
    items:
      type: string
    type: array
  secretNames:
    items:
      type: string
    type: array
required:
  - configMapNames
  - containerPort
  - deploymentNamespace
  - namespaces
  - secretNames
type: object
```

#### 5. Namespaces, ConfigMaps, Secrets, ContainerPort, and Resources (requested, limit) are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, ContainerPort, and Resources (requested, limit).

[config/crd/bases/resources-operator.kyma-project.io_companions.yaml](config/crd/bases/resources-operator.kyma-project.io_companions.yaml)

```yaml
properties:
  configMapNames:
    items:
      type: string
    type: array
  containerPort:
    format: int32
    type: integer
  deploymentNamespace:
    type: string
  namespaces:
    items:
      type: string
    type: array
  resources:
    description: ResourceRequirements defines the CPU and Memory requirements
      for the resources.
    properties:
      limits:
        description: ResourceValues defines the CPU and Memory values
          for the resources.
        properties:
          cpu:
            type: string
          memory:
            type: string
        required:
          - cpu
          - memory
        type: object
      requests:
        description: ResourceValues defines the CPU and Memory values
          for the resources.
        properties:
          cpu:
            type: string
          memory:
            type: string
        required:
          - cpu
          - memory
        type: object
    required:
      - limits
      - requests
    type: object
  secretNames:
    items:
      type: string
    type: array
required:
  - configMapNames
  - containerPort
  - deploymentNamespace
  - namespaces
  - resources
  - secretNames
type: object
```

#### 6. Namespaces, ConfigMaps, Secrets, ContainerPort, Resources (requested, limit), and Replicas are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, ContainerPort, Resources (requested, limit), and Replicas.

[config/crd/bases/replicas-operator.kyma-project.io_companions.yaml](config/crd/bases/replicas-operator.kyma-project.io_companions.yaml)

```yaml
properties:
  configMapNames:
    description: Required ConfigMaps names
    items:
      type: string
    type: array
  containerPort:
    description: Container port for the companion backend. Default value
      is 5000.
    format: int32
    type: integer
  deploymentNamespace:
    description: Namespace where the deployment will be created.
    type: string
  namespaces:
    description: |-
      List of namespaces which are prerequisites for the Kyma companion manager.
      Defaults:
      - 'ai-core': Main namespace. Namespace for the SAP AI Core component.
      - 'hana-cloud': Namespace for the SAP HANA Cloud vector instance.
      - 'redis': Namespace for the Redis.
    items:
      type: string
    type: array
  replicas:
    description: Replica count for the companion backend. Default value
      is 1.
    format: int32
    type: integer
  resources:
    description: |-
      Specify required resources and resource limits for the companion backend.
      Example:
      resources:
        limits:
          cpu: 1
          memory: 1Gi
        requests:
          cpu: 500m
          memory: 256Mi
    properties:
      limits:
        description: ResourceValues defines the CPU and Memory values
          for the resources.
        properties:
          cpu:
            type: string
          memory:
            type: string
        required:
          - cpu
          - memory
        type: object
      requests:
        description: ResourceValues defines the CPU and Memory values
          for the resources.
        properties:
          cpu:
            type: string
          memory:
            type: string
        required:
          - cpu
          - memory
        type: object
    required:
      - limits
      - requests
    type: object
  secretNames:
    description: Required Secrets names
    items:
      type: string
    type: array
required:
  - configMapNames
  - containerPort
  - deploymentNamespace
  - namespaces
  - replicas
  - resources
  - secretNames
type: object
```
