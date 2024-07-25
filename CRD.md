# Custom Resource Definitions (CRDs)

## Table of Contents

- [Companion Manager CR](#companion-manager-cr)
  - [Specification](#specification)
- [Deployment steps](#deployment-steps)
- [CRD Requirements](#crd-requirements)
- [CRD Options](#crd-options)
  - [1. All fields are hardcoded in the CRD.](#1-all-fields-are-hardcoded-in-the-crd)
  - [2. Namespaces are configurable.](#2-namespaces-are-configurable)
  - [3. Namespaces, ConfigMaps, and Secrets are configurable.](#3-namespaces-configmaps-and-secrets-are-configurable)
  - [4. Namespaces, ConfigMaps, Secrets and ContainerPort are configurable.](#4-namespaces-configmaps-secrets-and-containerport-are-configurable)
  - [5. Namespaces, ConfigMaps, Secrets, ContainerPort, and Resources (requested, limit) are configurable.](#5-namespaces-configmaps-secrets-containerport-and-resources-requested-limit-are-configurable)
  - [6. Namespaces, ConfigMaps, Secrets, ContainerPort, Resources (requested, limit), and Replicas are configurable.](#6-namespaces-configmaps-secrets-containerport-resources-requested-limit-and-replicas-are-configurable)
  - [7. Namespaces, ConfigMaps, Secrets, Resources (requested, limit), and Replicas are configurable.](#7-namespaces-configmaps-secrets-resources-requested-limit-and-replicas-are-configurable)
  - [8. Namespaces, ConfigMaps, Secrets, Resources (requested, limit), Replicas, Annotations, and Labels are configurable.](#8-namespaces-configmaps-secrets-resources-requested-limit-replicas-annotations-and-labels-are-configurable)
  - [9. Parameters are grouped by type (Companion, Hana, Redis)](#9-parameters-are-grouped-by-type-companion-hana-redis)
  - [10. Simple solution with secrets, resources and replicas (1)](#10-simple-solution-with-secrets-resources-and-replicas-1)
  - [11. Simple solution with secrets, resources, replicas (2)](#11-simple-solution-with-secrets-resources-replicas-2)
- [Conclusion - Suggestion for the CRD](#conclusion---suggestion-for-the-crd)
  - [Sample manifest 1](#sample-manifest-1)
  - [Sample manifest 2](#sample-manifest-2)
- [Other](#other)
  - [Default values](#default-values)
  - [Status.state](#statusstate)

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

## Deployment steps

This is a very high level overview of the deployment steps. This is not a detailed deployment guide.

![Deployment steps](docs/images/cr-deployment.drawio.png)

## CRD Requirements

From configuration perspective, not all fields are mandatory. The only mandatory fields depend on the following conditions:

- Flexibility
- Easy to deploy
- Easy to maintain
- Automation support

From this reason we should define the best Custom Resource Definition (CRD) for the Companion Manager CR.

## CRD Options

We could use the following CRD (Custom Resource Definition) options for the Companion Manager CR.

### 1. All fields are hardcoded in the CRD.

The user cannot change any field.

[config/crd/bases/1-operator.kyma-project.io_companions.yaml](config/crd/bases/1-operator.kyma-project.io_companions.yaml)

### 2. Namespaces are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace.

[config/crd/bases/2-operator.kyma-project.io_companions.yaml](config/crd/bases/2-operator.kyma-project.io_companions.yaml)

### 3. Namespaces, ConfigMaps, and Secrets are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, and Secrets.

[config/crd/bases/c3-operator.kyma-project.io_companions.yaml](config/crd/bases/3-operator.kyma-project.io_companions.yaml)

### 4. Namespaces, ConfigMaps, Secrets and ContainerPort are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, and ContainerPort.

[config/crd/bases/4-operator.kyma-project.io_companions.yaml](config/crd/bases/4-operator.kyma-project.io_companions.yaml)

### 5. Namespaces, ConfigMaps, Secrets, ContainerPort, and Resources (requested, limit) are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, ContainerPort, and Resources (requested, limit).

[config/crd/bases/5-operator.kyma-project.io_companions.yaml](config/crd/bases/5-operator.kyma-project.io_companions.yaml)

### 6. Namespaces, ConfigMaps, Secrets, ContainerPort, Resources (requested, limit), and Replicas are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, ContainerPort, Resources (requested, limit), and Replicas.

[config/crd/bases/6-operator.kyma-project.io_companions.yaml](config/crd/bases/6-operator.kyma-project.io_companions.yaml)

### 7. Namespaces, ConfigMaps, Secrets, Resources (requested, limit), and Replicas are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, Resources (requested, limit), and Replicas.

[config/crd/bases/7-operator.kyma-project.io_companions.yaml](config/crd/bases/7-operator.kyma-project.io_companions.yaml)

### 8. Namespaces, ConfigMaps, Secrets, Resources (requested, limit), Replicas, Annotations, and Labels are configurable.

The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace, ConfigMaps, Secrets, Resources (requested, limit), Replicas, Annotations, and Labels.

[config/crd/bases/8-operator.kyma-project.io_companions.yaml](config/crd/bases/8-operator.kyma-project.io_companions.yaml)

### 9. Parameters are grouped by type (Companion, Hana, Redis)

The user can change all important parameters for the companion application. The parameters are grouped by type (Companion, Hana, Redis).

[config/crd/bases/9-operator.kyma-project.io_companions.yaml](config/crd/bases/9-operator.kyma-project.io_companions.yaml)

### 10. Simple solution with secrets, resources and replicas (1)

The user can change the secrets, resources and replicas for the companion application.

Secrets are grouped by type (AI Core, Companion, Hana, Redis). One secrets is required for each type. Secret value contains the namespace and the secret name.
Resourses and replica settings are part of a companion section.

[config/crd/bases/10-operator.kyma-project.io_companions.yaml](config/crd/bases/10-operator.kyma-project.io_companions.yaml)

Sample manifest:
[config/samples/default-01.yaml](config/samples/default-01.yaml)

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: Companion
metadata:
  labels:
    app.kubernetes.io/name: default
    app.kubernetes.io/component: kyma-companion-manager
    app.kubernetes.io/part-of: kyma-companion-manager
  name: default
  namespace: kyma-system
spec:
  aiCoreSecret: ai-core/aicore
  companionSecret: companion/aicore
  hanaCloudSecret: companion/hana-cloud
  redisSecret: companion/redis
  companion:
    replicas:
      min: 1
      max: 3
    resources:
      limits:
        cpu: "4"
        memory: 4Gi
      requests:
        cpu: 500m
        memory: 256Mi
```

### 11. Simple solution with secrets, resources, replicas (2)

The user can change the secrets, resources and replicas for the companion application.

The related parameters are grouped by type (AI Core, Companion, Hana, Redis). In each group, the user can provide the related specific parameters. In our case, the user can provide the secret name and namespace for each type and the resources and replicas for the companion application.

[config/crd/bases/11-operator.kyma-project.io_companions.yaml](config/crd/bases/11-operator.kyma-project.io_companions.yaml)

Sample manifest:
[config/samples/default-02.yaml](config/samples/default-02.yaml)

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: Companion
metadata:
  labels:
    app.kubernetes.io/name: default
    app.kubernetes.io/component: kyma-companion-manager
    app.kubernetes.io/part-of: kyma-companion-manager
  name: default
  namespace: kyma-system
spec:
  aiCore:
    secret: ai-core/aicore
  hanaCloud:
    secret: companion/hana-cloud
  redis:
    secret: companion/redis
  companion:
    secret: companion/aicore
    replicas:
      min: 1
      max: 3
    resources:
      limits:
        cpu: "4"
        memory: 4Gi
      requests:
        cpu: 500m
        memory: 256Mi
```

## Conclusion - Suggestion for the CRD

The best option is the option 10-11. These provide the most flexibility for the user. The user can change the most important fields for the companion application, which ensures to be easy to deploy, easy to maintain, and automation support.


### Sample manifest 1

[config/samples/default.yaml](config/samples/default-01.yaml)

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: Companion
metadata:
  labels:
    app.kubernetes.io/name: default
    app.kubernetes.io/component: kyma-companion-manager
    app.kubernetes.io/part-of: kyma-companion-manager
  name: default
  namespace: kyma-system
spec:
  aiCoreSecret: ai-core/aicore
  companionSecret: companion/aicore
  hanaCloudSecret: companion/hana-cloud
  redisSecret: companion/redis
  companion:
    replicas:
      min: 1
      max: 3
    resources:
      limits:
        cpu: "4"
        memory: 4Gi
      requests:
        cpu: 500m
        memory: 256Mi
```

### Sample manifest 2

[config/samples/default.yaml](config/samples/default-02.yaml)

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: Companion
metadata:
  labels:
    app.kubernetes.io/name: default
    app.kubernetes.io/component: kyma-companion-manager
    app.kubernetes.io/part-of: kyma-companion-manager
  name: default
  namespace: kyma-system
spec:
  aiCore:
    secret: ai-core/aicore
  hanaCloud:
    secret: companion/hana-cloud
  redis:
    secret: companion/redis
  companion:
    secret: companion/aicore
    replicas:
      min: 1
      max: 3
    resources:
      limits:
        cpu: "4"
        memory: 4Gi
      requests:
        cpu: 500m
        memory: 256Mi
```

## Other

### Default values

In the final CRD, we should define the default values for the fields. The default values should be defined in the CRD.

### Status.state

Kyma modules should provide the `status.state`, because the Lifecycle Manager then updates the Kyma CR of the cluster based on the observed status changes in the module CR (similar to a native Kubernetes deployment tracking availability).

https://github.com/kyma-project/lifecycle-manager/blob/main/docs/modularization.md

```yaml
          status:
            description: CompanionStatus defines the observed state of Companion.
            properties:
              configMapsData:
                additionalProperties:
                  additionalProperties:
                    type: string
                  type: object
                description: 'ConfigMapsData: Map of ConfigMaps and their data. (optional)'
                type: object
              configMapsExists:
                additionalProperties:
                  type: boolean
                description: 'ConfigMapsExists: Map of ConfigMaps and their existence
                  status.'
                type: object
              namespacesExist:
                additionalProperties:
                  type: boolean
                description: |-
                  Result of prerequisites validation.
                  NamespacesExist: Map of namespaces and their existence status.
                type: object
              secretsData:
                additionalProperties:
                  additionalProperties:
                    format: byte
                    type: string
                  type: object
                description: 'SecretsData: Map of Secrets and their data. (optional)'
                type: object
              secretsExists:
                additionalProperties:
                  type: boolean
                description: 'SecretsExists: Map of Secrets and their existence status.'
                type: object
              state:
                description: |-
                  Defines the overall state of the Companion custom resource.<br/>
                  - `Ready` when all the resources managed by the Kyma companion manager are deployed successfully and
                  the companion backend is ready.<br/>
                  - `Warning` if there is a user input misconfiguration.<br/>
                  - `Processing` if the resources managed by the Kyma companion manager are being created or updated.<br/>
                  - `Error` if an error occurred while reconciling the Companion custom resource.
                  - `Deleting` if the resources managed by the Kyma companion manager are being deleted.
                type: string
            required:
            - configMapsExists
            - namespacesExist
            - secretsExists
            - state
            type: object
        type: object
```
