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

#### 1. All fields are hardcoded in the CRD. The user cannot change any field.

[config/crd/bases/all-hardcoded-operator.kyma-project.io_companions.yaml](config/crd/bases/all-hardcoded-operator.kyma-project.io_companions.yaml)

````yaml
            properties:
              state:
                description: |-
                  Defines the overall state of the Companion custom resource.<br/>
                  - `Ready` when all the resources managed by the Kyma companion manager are deployed successfully and
                  the companion backend is ready.<br/>
                  - `Warning` if there is a user input misconfiguration.<br/>
                  - `Processing` if the resources managed by the Kyma companion manager are being created or updated.<br/>
                  - `Error` if an error occurred while reconciling the Companion custom resource.
                type: string
            required:
            - state
            type: object
            
````


#### 2. Namespaces are configurable. The user can change the namespaces in which are required for the companion application. (SAP AI Core, Redis, Hana Cloud) User also can change the deployment namespace.

[config/crd/bases/namespaces-operator.kyma-project.io_companions.yaml](config/crd/bases/namespaces-operator.kyma-project.io_companions.yaml)

````yaml
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
          status:
            description: CompanionStatus defines the observed state of Companion.
            properties:
              state:
                description: |-
                  Defines the overall state of the Companion custom resource.<br/>
                  - `Ready` when all the resources managed by the Kyma companion manager are deployed successfully and
                  the companion backend is ready.<br/>
                  - `Warning` if there is a user input misconfiguration.<br/>
                  - `Processing` if the resources managed by the Kyma companion manager are being created or updated.<br/>
                  - `Error` if an error occurred while reconciling the Companion custom resource.
                type: string
            required:
            - state
            type: object
````