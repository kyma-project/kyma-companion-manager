package controller

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"github.com/stretchr/testify/require"
	kadmissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// MockedUnitTestEnvironment provides mocked resources for unit tests.
type MockedUnitTestEnvironment struct {
	Client     client.Client
	Reconciler *Reconciler
	Recorder   *record.FakeRecorder
}

func NewMockedUnitTestEnvironment(t *testing.T, objs ...client.Object) *MockedUnitTestEnvironment {
	t.Helper()

	// setup fake client for k8s
	newScheme := runtime.NewScheme()
	err := kcmv1alpha1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kcorev1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kadmissionregistrationv1.AddToScheme(newScheme)
	require.NoError(t, err)

	// Create k8s client.
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(newScheme)
	fakeClient := fakeClientBuilder.WithObjects(objs...).WithStatusSubresource(objs...).Build()

	// fake recorder.
	recorder := &record.FakeRecorder{}

	// setup reconciler
	reconciler := &Reconciler{
		Client: fakeClient,
		Scheme: newScheme,
	}

	return &MockedUnitTestEnvironment{
		Client:     fakeClient,
		Reconciler: reconciler,
		Recorder:   recorder,
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
