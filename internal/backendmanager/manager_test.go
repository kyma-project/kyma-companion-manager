package backendmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kcmlabel "github.com/kyma-project/kyma-companion-manager/internal/label"
	kcmk8smocks "github.com/kyma-project/kyma-companion-manager/pkg/k8s/mocks"
	"github.com/kyma-project/kyma-companion-manager/pkg/utils"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

func Test_GenerateNewDeployment(t *testing.T) {
	t.Parallel()

	// given
	givenBackendImage := "kyma-project/backend:11072024"
	givenCompanion := testutils.NewCompanionCR()
	logger, err := testutils.NewSugaredLogger()
	givenTerminationGracePeriodSeconds := terminationGracePeriodSeconds
	require.NoError(t, err)
	backendManager := NewBackendManager(nil, nil, logger)

	// when
	gotDeployment, err := backendManager.GenerateNewDeployment(givenCompanion, givenBackendImage)

	// then
	require.NoError(t, err)
	wantDeployment := &kappsv1.Deployment{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:            BackendResourceName,
			Namespace:       givenCompanion.Namespace,
			Labels:          kcmlabel.GetCommonLabels(BackendResourceName),
			OwnerReferences: getOwnerReferences(*givenCompanion),
		},
		Spec: kappsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(deploymentReplicas),
			Selector: kmetav1.SetAsLabelSelector(kcmlabel.GetCommonLabels(BackendResourceName)),
			Template: kcorev1.PodTemplateSpec{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:   BackendResourceName,
					Labels: kcmlabel.GetCommonLabels(BackendResourceName),
				},
				Spec: kcorev1.PodSpec{
					RestartPolicy: kcorev1.RestartPolicyAlways,
					// SecurityContext:               getPodSecurityContext(),
					TerminationGracePeriodSeconds: &givenTerminationGracePeriodSeconds,
					PriorityClassName:             priorityClassName,
					Containers: []kcorev1.Container{
						{
							Name:            kcmlabel.ValueCompanionBackend,
							Image:           givenBackendImage,
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
					},
					Volumes: []kcorev1.Volume{
						{
							Name: BackendResourceName,
							VolumeSource: kcorev1.VolumeSource{
								Secret: &kcorev1.SecretVolumeSource{
									SecretName:  BackendResourceName,
									DefaultMode: utils.Int32Ptr(420),
								},
							},
						},
					},
				},
			},
		},
		Status: kappsv1.DeploymentStatus{},
	}
	// compare pointers of bool.
	require.Len(t, gotDeployment.OwnerReferences, 1)
	require.Equal(t, *wantDeployment.OwnerReferences[0].BlockOwnerDeletion,
		*gotDeployment.OwnerReferences[0].BlockOwnerDeletion)
	require.Equal(t, *wantDeployment.OwnerReferences[0].Controller, *gotDeployment.OwnerReferences[0].Controller)

	gotDeployment.OwnerReferences[0].BlockOwnerDeletion = wantDeployment.OwnerReferences[0].BlockOwnerDeletion
	gotDeployment.OwnerReferences[0].Controller = wantDeployment.OwnerReferences[0].Controller

	// compare object.
	require.Equal(t, wantDeployment, gotDeployment)
}

func Test_GenerateNewSecret(t *testing.T) {
	t.Parallel()

	// given
	givenCompanion := testutils.NewCompanionCR()
	logger, err := testutils.NewSugaredLogger()
	require.NoError(t, err)
	backendManager := NewBackendManager(nil, nil, logger)

	givenConfig := Config{
		HanaDB:       []byte("hanaDB"),
		Redis:        []byte("redis"),
		AICoreSecret: []byte("ai-core-secret"),
		AICoreConfig: []byte("ai-core-config"),
	}

	// when
	gotSecret, err := backendManager.GenerateNewSecret(givenCompanion, givenConfig)

	// then
	require.NoError(t, err)
	wantSecret := &kcorev1.Secret{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:            BackendResourceName,
			Namespace:       givenCompanion.Namespace,
			Labels:          kcmlabel.GetCommonLabels(BackendResourceName),
			OwnerReferences: getOwnerReferences(*givenCompanion),
		},
		Data: map[string][]byte{
			"hana-db-secret": givenConfig.HanaDB,
			"redis-secret":   givenConfig.Redis,
			"ai-core-config": givenConfig.AICoreConfig,
			"ai-core-secret": givenConfig.AICoreSecret,
		},
		Type: kcorev1.SecretTypeOpaque,
	}

	// compare pointers of bool.
	require.Len(t, gotSecret.OwnerReferences, 1)
	require.Equal(t, *wantSecret.OwnerReferences[0].BlockOwnerDeletion,
		*gotSecret.OwnerReferences[0].BlockOwnerDeletion)
	require.Equal(t, *wantSecret.OwnerReferences[0].Controller, *gotSecret.OwnerReferences[0].Controller)

	gotSecret.OwnerReferences[0].BlockOwnerDeletion = wantSecret.OwnerReferences[0].BlockOwnerDeletion
	gotSecret.OwnerReferences[0].Controller = wantSecret.OwnerReferences[0].Controller

	// compare object.
	require.Equal(t, wantSecret, gotSecret)
}

func Test_GetBackendConfig(t *testing.T) {
	t.Parallel()

	// given
	logger, err := testutils.NewSugaredLogger()
	require.NoError(t, err)

	kubeClient := new(kcmk8smocks.Client)
	backendManager := NewBackendManager(nil, kubeClient, logger)

	// define mock behaviour.
	sampleSecret := &kcorev1.Secret{
		Data: map[string][]byte{
			"test-key": []byte("test-value"),
		},
	}
	sampleConfigMap := &kcorev1.ConfigMap{
		Data: map[string]string{
			"test-key": "test-value",
		},
	}
	kubeClient.On("GetSecret", mock.Anything, mock.Anything, mock.Anything).Return(
		sampleSecret, nil).Times(3)
	kubeClient.On("GetConfigMap", mock.Anything, mock.Anything, mock.Anything).Return(
		sampleConfigMap, nil).Once()

	// when
	gotConfig, err := backendManager.GetBackendConfig(context.TODO())

	// then
	require.NoError(t, err)
	sampleSecretData := "{\"test-key\":\"dGVzdC12YWx1ZQ==\"}"
	sampleConfigData := "{\"test-key\":\"test-value\"}"

	wantConfig := &Config{
		HanaDB:       []byte(sampleSecretData),
		Redis:        []byte(sampleSecretData),
		AICoreSecret: []byte(sampleSecretData),
		AICoreConfig: []byte(sampleConfigData),
	}

	// compare object.
	require.Equal(t, wantConfig, gotConfig)
	kubeClient.AssertExpectations(t)
}
