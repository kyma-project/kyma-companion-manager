package controller

// **NOTE:** This file contains unit tests for companion_controller.go.
// Integration tests for controller are located in the test/integration directory.

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"github.com/kyma-project/kyma-companion-manager/internal/backendmanager"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

func Test_loggerWithCompanion(t *testing.T) {
	t.Parallel()

	// given
	testEnv := NewMockedUnitTestEnvironment(t)
	givenCompanion := testutils.NewCompanionCR()

	// when
	gotLogger := testEnv.Reconciler.loggerWithCompanion(givenCompanion)

	// then
	require.NotNil(t, gotLogger)
}

func Test_reconcileDeployment(t *testing.T) {
	t.Parallel()

	// define test cases
	testCases := []struct {
		name                    string
		givenCompanion          *kcmv1alpha1.Companion
		givenMocksBehaviourFunc func(testEnv *MockedUnitTestEnvironment, deployment *kappsv1.Deployment)
	}{
		{
			name:           "should update the deployment when it does not exist",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenDeployment *kappsv1.Deployment) {
				testEnv.backendManager.On("GenerateNewDeployment",
					mock.Anything, mock.Anything).Return(givenDeployment, nil).Once()
				testEnv.kubeClient.On("GetDeployment",
					mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
				testEnv.kubeClient.On("PatchApply",
					mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name:           "should not update the deployment when it exists",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenDeployment *kappsv1.Deployment) {
				testEnv.backendManager.On("GenerateNewDeployment",
					mock.Anything, mock.Anything).Return(givenDeployment, nil).Once()
				testEnv.kubeClient.On("GetDeployment",
					mock.Anything, mock.Anything, mock.Anything).Return(givenDeployment, nil).Once()
			},
		},
		{
			name:           "should update the deployment when the existing deployment is different from expected",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenDeployment *kappsv1.Deployment) {
				testEnv.backendManager.On("GenerateNewDeployment",
					mock.Anything, mock.Anything).Return(givenDeployment, nil).Once()

				changedDeployment := givenDeployment.DeepCopy()
				changedDeployment.Spec.Template.Spec.Containers[0].Image = "changed-image"
				testEnv.kubeClient.On("GetDeployment",
					mock.Anything, mock.Anything, mock.Anything).Return(changedDeployment, nil).Once()
				testEnv.kubeClient.On("PatchApply",
					mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// given
			givenDeployment := testutils.NewCompanionDeployment("test-deployment",
				"test-namespace")
			testEnv := NewMockedUnitTestEnvironment(t, tc.givenCompanion)

			// define mocks behaviour
			tc.givenMocksBehaviourFunc(testEnv, givenDeployment)

			// when
			err := testEnv.Reconciler.reconcileDeployment(context.TODO(), tc.givenCompanion, testEnv.Logger)

			// then
			require.NoError(t, err)
			testEnv.backendManager.AssertExpectations(t)
			testEnv.kubeClient.AssertExpectations(t)
		})
	}
}

func Test_reconcileSecret(t *testing.T) {
	t.Parallel()

	// define test cases
	testCases := []struct {
		name                    string
		givenCompanion          *kcmv1alpha1.Companion
		givenMocksBehaviourFunc func(testEnv *MockedUnitTestEnvironment, givenSecret *kcorev1.Secret,
			givenConfig *backendmanager.Config)
	}{
		{
			name:           "should update the secret when it does not exist",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenSecret *kcorev1.Secret,
				givenConfig *backendmanager.Config,
			) {
				testEnv.backendManager.On("GenerateNewSecret",
					mock.Anything, mock.Anything).Return(givenSecret, nil).Once()
				testEnv.backendManager.On("GetBackendConfig",
					mock.Anything, mock.Anything).Return(givenConfig, nil).Once()
				testEnv.kubeClient.On("GetSecret",
					mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
				testEnv.kubeClient.On("PatchApply",
					mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name:           "should not update the secret when it exists",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenSecret *kcorev1.Secret,
				givenConfig *backendmanager.Config,
			) {
				testEnv.backendManager.On("GenerateNewSecret",
					mock.Anything, mock.Anything).Return(givenSecret, nil).Once()
				testEnv.backendManager.On("GetBackendConfig",
					mock.Anything, mock.Anything).Return(givenConfig, nil).Once()
				testEnv.kubeClient.On("GetSecret",
					mock.Anything, mock.Anything, mock.Anything).Return(givenSecret, nil).Once()
			},
		},
		{
			name:           "should update the secret when the existing secret is different from expected",
			givenCompanion: testutils.NewCompanionCR(),
			givenMocksBehaviourFunc: func(testEnv *MockedUnitTestEnvironment, givenSecret *kcorev1.Secret,
				givenConfig *backendmanager.Config,
			) {
				testEnv.backendManager.On("GenerateNewSecret",
					mock.Anything, mock.Anything).Return(givenSecret, nil).Once()
				testEnv.backendManager.On("GetBackendConfig",
					mock.Anything, mock.Anything).Return(givenConfig, nil).Once()

				changedSecret := givenSecret.DeepCopy()
				changedSecret.Data = map[string][]byte{
					"test-key": []byte("changed-value"),
				}
				testEnv.kubeClient.On("GetSecret",
					mock.Anything, mock.Anything, mock.Anything).Return(changedSecret, nil).Once()
				testEnv.kubeClient.On("PatchApply",
					mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// given
			givenSecret := &kcorev1.Secret{
				Data: map[string][]byte{
					"test-key": []byte("test-value"),
				},
			}
			testEnv := NewMockedUnitTestEnvironment(t, tc.givenCompanion)

			// define mocks behaviour
			tc.givenMocksBehaviourFunc(testEnv, givenSecret, &backendmanager.Config{})

			// when
			err := testEnv.Reconciler.reconcileSecret(context.TODO(), tc.givenCompanion, testEnv.Logger)

			// then
			require.NoError(t, err)
			testEnv.backendManager.AssertExpectations(t)
			testEnv.kubeClient.AssertExpectations(t)
		})
	}
}
