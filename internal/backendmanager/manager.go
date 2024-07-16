package backendmanager

import (
	"context"
	"encoding/json"
	"strings"

	"go.uber.org/zap"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	kcmlabel "github.com/kyma-project/kyma-companion-manager/internal/label"
	kcmk8s "github.com/kyma-project/kyma-companion-manager/pkg/k8s"
	kcmk8sdeployment "github.com/kyma-project/kyma-companion-manager/pkg/k8s/deployment"
	kcmk8ssecret "github.com/kyma-project/kyma-companion-manager/pkg/k8s/secret"
)

const (
	priorityClassName             = "kyma-companion-manager-priority-class"
	backendPortName               = "http"
	backendPortNum                = int32(8000)
	backendMetricsPortName        = "http-metrics"
	backendMetricsPortNum         = int32(9090)
	livenessInitialDelaySecs      = int32(5)
	livenessTimeoutSecs           = int32(1)
	livenessPeriodSecs            = int32(2)
	terminationGracePeriodSeconds = int64(30)
	requestsCPU                   = "200m"
	requestsMemory                = "512Mi"
	limitsCPU                     = "500m"
	limitsMemory                  = "1Gi"
	deploymentReplicas            = 1
	secretMountPath               = "/mnt/secrets"
)

// compile-time check.
var _ Manager = &BackendManager{}

//go:generate go run github.com/vektra/mockery/v2 --name=Manager --outpkg=mocks --case=underscore
type Manager interface {
	GenerateNewDeployment(companion *kcmv1alpha1.Companion, backendImage string) (*kappsv1.Deployment, error)
	GenerateNewSecret(companion *kcmv1alpha1.Companion, config Config) (*kcorev1.Secret, error)
	GetBackendConfig(ctx context.Context) (*Config, error)
}

type BackendManager struct {
	client.Client
	kubeClient kcmk8s.Client
	logger     *zap.SugaredLogger
}

func NewBackendManager(
	client client.Client,
	kubeClient kcmk8s.Client,
	logger *zap.SugaredLogger,
) Manager {
	return &BackendManager{
		Client:     client,
		kubeClient: kubeClient,
		logger:     logger,
	}
}

func (m *BackendManager) GenerateNewDeployment(companion *kcmv1alpha1.Companion,
	backendImage string,
) (*kappsv1.Deployment, error) {
	// define labels.
	labels := kcmlabel.GetCommonLabels(BackendResourceName)

	// define containers.
	containers := []kcorev1.Container{
		{
			Name:            kcmlabel.ValueCompanionBackend,
			Image:           backendImage,
			Ports:           getContainerPorts(),
			LivenessProbe:   getLivenessProbe(),
			ReadinessProbe:  getReadinessProbe(),
			ImagePullPolicy: kcorev1.PullAlways,
			// SecurityContext: getContainerSecurityContext(),
			Resources: getResources(requestsCPU, requestsMemory, limitsCPU, limitsMemory),
			VolumeMounts: []kcorev1.VolumeMount{
				{
					Name:      BackendResourceName,
					ReadOnly:  true,
					MountPath: secretMountPath,
				},
			},
		},
	}

	// define deployment object.
	deployment := kcmk8sdeployment.NewDeployment(
		BackendResourceName,
		companion.GetNamespace(),
		kcmk8sdeployment.WithLabels(labels),
		kcmk8sdeployment.WithRestartPolicyAlways(),
		kcmk8sdeployment.WithReplicas(deploymentReplicas),
		// kcmk8sdeployment.WithSecurityContext(getPodSecurityContext()),
		kcmk8sdeployment.WithTerminationGracePeriodSeconds(terminationGracePeriodSeconds),
		kcmk8sdeployment.WithPriorityClassName(priorityClassName),
		kcmk8sdeployment.WithSelectorLabels(labels),
		kcmk8sdeployment.WithContainers(containers),
		kcmk8sdeployment.WithOwnerReferences(getOwnerReferences(*companion)),
		kcmk8sdeployment.WithVolumeMountedSecret(BackendResourceName),
	)

	return deployment, nil
}

func (m *BackendManager) GenerateNewSecret(companion *kcmv1alpha1.Companion, config Config) (*kcorev1.Secret, error) {
	// define secret object.
	secret := kcmk8ssecret.NewSecret(
		BackendResourceName,
		companion.GetNamespace(),
		kcmk8ssecret.WithLabels(kcmlabel.GetCommonLabels(BackendResourceName)),
		kcmk8ssecret.WithOwnerReferences(getOwnerReferences(*companion)),
		kcmk8ssecret.WithDataKeyKey("hana-db-secret", config.HanaDB),
		kcmk8ssecret.WithDataKeyKey("redis-secret", config.Redis),
		kcmk8ssecret.WithDataKeyKey("ai-core-config", config.AICoreConfig),
		kcmk8ssecret.WithDataKeyKey("ai-core-secret", config.AICoreSecret),
	)

	return secret, nil
}

func (m *BackendManager) GetBackendConfig(ctx context.Context) (*Config, error) {
	// define config object.
	config := &Config{}

	// hard-coded secret namespaced names for first iteration.
	hanaDBSecretNamespacedName := "kyma-system/companion-hana-db"
	redisSecretNamespacedName := "kyma-system/companion-redis"
	aiCoreSecretNamespacedName := "kyma-system/companion-ai-core"
	aiCoreConfigNamespacedName := "kyma-system/companion-ai-core"

	// Fetch the secret for HANA Vector DB.
	parts := strings.Split(hanaDBSecretNamespacedName, "/")
	hanaDBSecret, err := m.kubeClient.GetSecret(ctx, parts[1], parts[0])
	if err != nil {
		return nil, err
	}
	jsonString, err := json.Marshal(hanaDBSecret.Data)
	if err != nil {
		return nil, err
	}
	config.HanaDB = jsonString

	// Fetch the secret for Redis.
	parts = strings.Split(redisSecretNamespacedName, "/")
	redisSecret, err := m.kubeClient.GetSecret(ctx, parts[1], parts[0])
	if err != nil {
		return nil, err
	}
	jsonString, err = json.Marshal(redisSecret.Data)
	if err != nil {
		return nil, err
	}
	config.Redis = jsonString

	// Fetch the secret for AI-Core.
	parts = strings.Split(aiCoreSecretNamespacedName, "/")
	aiCoreSecret, err := m.kubeClient.GetSecret(ctx, parts[1], parts[0])
	if err != nil {
		return nil, err
	}
	jsonString, err = json.Marshal(aiCoreSecret.Data)
	if err != nil {
		return nil, err
	}
	config.AICoreSecret = jsonString

	// Fetch the configMap for AI-Core.
	parts = strings.Split(aiCoreConfigNamespacedName, "/")
	aiCoreConfigMap, err := m.kubeClient.GetConfigMap(ctx, parts[1], parts[0])
	if err != nil {
		return nil, err
	}
	jsonString, err = json.Marshal(aiCoreConfigMap.Data)
	if err != nil {
		return nil, err
	}
	config.AICoreConfig = jsonString

	return config, nil
}
