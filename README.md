# Kyma Companion Manager
<!--- mandatory --->

[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/kyma-companion-manager)](https://api.reuse.software/info/github.com/kyma-project/kyma-companion-manager)

## Overview
<!--- mandatory section --->

Kyma Companion Manager is a standard Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that observes the state of companion resources and reconciles their state according to the desired state. It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/), which provide a reconcile function responsible for synchronizing resources until the desired state is reached in the cluster.

This project is scaffolded using [Kubebuilder](https://book.kubebuilder.io), and all the Kubebuilder `makefile` helpers mentioned [here](https://book.kubebuilder.io/reference/makefile-helpers.html) can be used.

## Get Started

You need a Kubernetes cluster to run against. You can use [k3d](https://k3d.io/) to get a local cluster for testing, or run against a remote cluster.
> [!NOTE]
> Your controller automatically uses the current context in your kubeconfig file, that is, whatever cluster `kubectl cluster-info` shows.

## Development

### Prerequisites

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kubebuilder](https://book.kubebuilder.io/)
- [kustomize](https://kustomize.io/)
- Access to Kubernetes cluster ([k3d](https://k3d.io/) / k8s)

### Run Kyma Companion Manager Locally

1. Install the CRDs into the cluster:

   ```sh
   make install
   ```

2. Run Kyma Companion Manager. It runs in the foreground, so if you want to leave it running, switch to a new terminal.

   ```sh
   make run
   ```

> [!NOTE]
> You can also run this in one step with the command: `make install run`.

### Run Tests

Run the unit and integration tests:

```sh
make generate-and-test
```

### Linting

1. Fix common lint issues:

   ```sh
   make imports
   make fmt
   make lint
   ```

### Modify the API Definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs:

```sh
make manifests
```

> [!NOTE]
> Run `make --help` for more information on all potential `make` targets.

For more information, see the [Kubebuilder documentation](https://book.kubebuilder.io/introduction.html).

### Build Container Images

Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<container-registry>/kyma-companion-manager:<tag> # If using docker, <container-registry> is your username.
```

> **NOTE**: For MacBook M1 devices, run:

```sh
make docker-buildx IMG=<container-registry>/kyma-companion-manager:<tag>
```
## Deployment

You need a Kubernetes cluster to run against. You can use [k3d](https://k3d.io/) to get a local cluster for testing, or run against a remote cluster.
> [!NOTE]
> Your controller automatically uses the current context in your kubeconfig file, that is, whatever cluster `kubectl cluster-info` shows.

### Deploy in the Cluster

1. Download Go packages:

   ```sh
   go mod vendor && go mod tidy
   ```

2. Install the CRDs to the cluster:

   ```sh
   make install
   ```

3. Build and push your image to the location specified by `IMG`:

   ```sh
   make docker-build docker-push IMG=<container-registry>/kyma-companion-manager:<tag>
   ```

4. Deploy the `kyma-companion-manager` controller to the cluster:

   ```sh
   make deploy IMG=<container-registry>/kyma-companion-manager:<tag>
   ```

5. [Optional] Install `Companion` Custom Resource:

   ```sh
   kubectl apply -f config/samples/default.yaml
   ```

### Undeploy Kyma Companion Manager

Undeploy Kyma Companion Manager from the cluster:

   ```sh
   make undeploy
   ```

### Uninstall CRDs

To delete the CRDs from the cluster:

   ```sh
   make uninstall
   ```

## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md).

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
