package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kapierrors "k8s.io/apimachinery/pkg/api/errors"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

func Test_UpdateDeployment(t *testing.T) {
	t.Parallel()

	// Define test cases
	testCases := []struct {
		name                   string
		namespace              string
		givenNewDeploymentSpec kappsv1.DeploymentSpec
		givenDeploymentExists  bool
	}{
		{
			name:                  "should update the deployment",
			namespace:             "test-namespace-1",
			givenDeploymentExists: true,
		},
		{
			name:                  "should give error that deployment does not exist",
			namespace:             "test-namespace-2",
			givenDeploymentExists: false,
		},
	}

	// Run tests
	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			// given
			ctx := context.Background()
			fakeClient := fake.NewClientBuilder().Build()
			kubeClient := &KubeClient{
				client: fakeClient,
			}
			givenDeployment := testutils.NewDeployment("test-deployment", testcase.namespace, map[string]string{})
			// Create the deployment if it should exist
			if testcase.givenDeploymentExists {
				require.NoError(t, fakeClient.Create(ctx, givenDeployment))
			}

			givenUpdatedDeployment := givenDeployment.DeepCopy()
			givenUpdatedDeployment.Spec = testcase.givenNewDeploymentSpec

			// when
			err := kubeClient.UpdateDeployment(ctx, givenUpdatedDeployment)

			// then
			if !testcase.givenDeploymentExists {
				require.Error(t, err)
				require.True(t, kapierrors.IsNotFound(err))
			} else {
				gotDeploy, err := kubeClient.GetDeployment(ctx, givenDeployment.Name, givenDeployment.Namespace)
				require.NoError(t, err)
				require.Equal(t, testcase.givenNewDeploymentSpec, gotDeploy.Spec)
			}
		})
	}
}

func Test_DeleteResource(t *testing.T) {
	t.Parallel()
	// Define test cases
	testCases := []struct {
		name                  string
		givenDeployment       *kappsv1.Deployment
		givenDeploymentExists bool
	}{
		{
			name: "should delete the deployment",
			givenDeployment: &kappsv1.Deployment{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
			},
			givenDeploymentExists: true,
		},
		{
			name: "should not return error when the deployment does not exist",
			givenDeployment: &kappsv1.Deployment{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace2",
				},
			},
			givenDeploymentExists: false,
		},
	}

	// Run tests
	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			// given
			ctx := context.Background()
			var givenObjs []client.Object
			if testcase.givenDeploymentExists {
				givenObjs = append(givenObjs, testcase.givenDeployment)
			}
			fakeClient := fake.NewClientBuilder().WithObjects(givenObjs...).Build()
			kubeClient := &KubeClient{
				client: fakeClient,
			}

			// when
			err := kubeClient.DeleteResource(ctx, testcase.givenDeployment)

			// then
			require.NoError(t, err)
			// Check that the deployment must not exist.
			err = fakeClient.Get(ctx, types.NamespacedName{
				Name:      testcase.givenDeployment.Name,
				Namespace: testcase.givenDeployment.Namespace,
			}, &kappsv1.Deployment{})
			require.True(t, kapierrors.IsNotFound(err), "DeleteResource did not delete deployment")
		})
	}
}

func Test_DeleteDeployment(t *testing.T) {
	t.Parallel()
	// Define test cases
	testCases := []struct {
		name         string
		namespace    string
		noDeployment bool
	}{
		{
			name:      "deployment exists",
			namespace: "test-namespace",
		},
		{
			name:         "deployment does not exist",
			namespace:    "test-namespace",
			noDeployment: true,
		},
	}

	// Run tests
	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			// given
			ctx := context.Background()
			fakeClient := fake.NewClientBuilder().Build()
			kubeClient := &KubeClient{
				client: fakeClient,
			}
			deployment := &kappsv1.Deployment{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
			}
			// Create the deployment if it should exist
			if !testcase.noDeployment {
				if err := fakeClient.Create(ctx, deployment); err != nil {
					t.Fatalf("failed to create deployment: %v", err)
				}
			}

			// when
			err := kubeClient.DeleteDeployment(ctx, deployment.Name, deployment.Namespace)

			// then
			require.NoError(t, err)
			// Check that the deployment was deleted
			err = fakeClient.Get(ctx,
				types.NamespacedName{Name: "test-deployment", Namespace: testcase.namespace}, &kappsv1.Deployment{})
			require.True(t, kapierrors.IsNotFound(err), "DeleteDeployment did not delete deployment")
		})
	}
}

func Test_GetSecret(t *testing.T) {
	t.Parallel()
	// Define test cases as a table.
	testCases := []struct {
		name              string
		givenName         string
		givenNamespace    string
		wantSecret        *kcorev1.Secret
		wantError         error
		wantNotFoundError bool
	}{
		{
			name:           "success",
			givenName:      "test-secret",
			givenNamespace: "test-namespace",
			wantSecret: &kcorev1.Secret{
				TypeMeta: kmetav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "test-namespace",
				},
				Data: map[string][]byte{
					"key": []byte("value"),
				},
			},
		},
		{
			name:              "not found",
			givenName:         "test-secret",
			givenNamespace:    "test-namespace",
			wantSecret:        nil,
			wantNotFoundError: true,
		},
	}

	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			// given
			ctx := context.Background()
			fakeClient := fake.NewClientBuilder().Build()
			kubeClient := &KubeClient{
				client: fakeClient,
			}

			// Create the secret if it should exist
			if testcase.wantSecret != nil {
				require.NoError(t, fakeClient.Create(ctx, testcase.wantSecret))
			}

			// Call the GetSecret function with the test case's givenNamespacedName.
			secret, err := kubeClient.GetSecret(context.Background(), testcase.givenName, testcase.givenNamespace)

			// Assert that the function returned the expected secret and error.
			if testcase.wantNotFoundError {
				require.True(t, kapierrors.IsNotFound(err))
			} else {
				require.ErrorIs(t, err, testcase.wantError)
			}
			require.Equal(t, testcase.wantSecret, secret)
		})
	}
}

func Test_GetConfigMap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		givenName         string
		givenNamespace    string
		wantNotFoundError bool
	}{
		{
			name:              "should return configmap",
			givenName:         "test-name",
			givenNamespace:    "test-namespace",
			wantNotFoundError: false,
		},
		{
			name:              "should not return configmap",
			givenName:         "non-existing",
			givenNamespace:    "non-existing",
			wantNotFoundError: true,
		},
	}

	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			// given
			ctx := context.Background()
			kubeClient := &KubeClient{client: fake.NewClientBuilder().Build()}
			givenCM := testutils.NewConfigMap(testcase.givenName, testcase.givenNamespace)
			if !testcase.wantNotFoundError {
				require.NoError(t, kubeClient.client.Create(ctx, givenCM))
			}

			// when
			gotCM, err := kubeClient.GetConfigMap(context.Background(), testcase.givenName, testcase.givenNamespace)

			// then
			if testcase.wantNotFoundError {
				require.Error(t, err)
				require.True(t, kapierrors.IsNotFound(err))
			} else {
				require.NoError(t, err)
				require.Equal(t, givenCM.GetName(), gotCM.Name)
			}
		})
	}
}
