package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	kadmissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	backendmanagermocks "github.com/kyma-project/kyma-companion-manager/internal/backendmanager/mocks"
	kcmk8smocks "github.com/kyma-project/kyma-companion-manager/pkg/k8s/mocks"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

// MockedUnitTestEnvironment provides mocked resources for unit tests.
type MockedUnitTestEnvironment struct {
	Client         client.Client
	Reconciler     *Reconciler
	Recorder       *record.FakeRecorder
	Logger         *zap.SugaredLogger
	kubeClient     *kcmk8smocks.Client
	backendManager *backendmanagermocks.Manager
}

func NewMockedUnitTestEnvironment(t *testing.T, objs ...client.Object) *MockedUnitTestEnvironment {
	t.Helper()

	// setup logger
	logger, err := testutils.NewSugaredLogger()
	require.NoError(t, err)

	// setup fake client for k8s
	newScheme := runtime.NewScheme()
	err = kcmv1alpha1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kcorev1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kadmissionregistrationv1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kappsv1.AddToScheme(newScheme)
	require.NoError(t, err)

	// Create k8s client.
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(newScheme)
	fakeClient := fakeClientBuilder.WithObjects(objs...).WithStatusSubresource(objs...).Build()

	// fake recorder.
	recorder := &record.FakeRecorder{}

	// setup custom mocks
	backendManager := new(backendmanagermocks.Manager)
	kubeClient := new(kcmk8smocks.Client)

	// setup reconciler
	reconciler := &Reconciler{
		Client:         fakeClient,
		Scheme:         newScheme,
		logger:         logger,
		kubeClient:     kubeClient,
		backendManager: backendManager,
	}

	return &MockedUnitTestEnvironment{
		Client:         fakeClient,
		Reconciler:     reconciler,
		Recorder:       recorder,
		Logger:         logger,
		kubeClient:     kubeClient,
		backendManager: backendManager,
	}
}

func (testEnv *MockedUnitTestEnvironment) GetCompanion(name, namespace string) (kcmv1alpha1.Companion, error) {
	var companion kcmv1alpha1.Companion
	err := testEnv.Client.Get(context.Background(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &companion)
	return companion, err
}
