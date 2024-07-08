package controller_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"github.com/kyma-project/kyma-companion-manager/test/integration"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

const projectRootDir = "../../../"

var testEnvironment *integration.TestEnvironment

// TestMain pre-hook and post-hook to run before and after all tests.
func TestMain(m *testing.M) {
	// Note: The setup will provision a single K8s env and
	// all the tests need to create and use a separate namespace

	// setup env test
	var err error
	testEnvironment, err = integration.NewTestEnvironment(projectRootDir, nil)
	if err != nil {
		log.Fatal(err)
	}

	// run tests
	code := m.Run()

	// tear down test env
	if err = testEnvironment.TearDown(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func Test_CreateCompanionCR(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		givenCompanion *kcmv1alpha1.Companion
	}{
		{
			name:           "dummy-test-case Companion CR should have ready status when deployment is ready",
			givenCompanion: testutils.NewCompanionCR(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// given
			// create unique namespace for this test run.
			givenNamespace := tc.givenCompanion.GetNamespace()
			testEnvironment.EnsureNamespaceCreation(t, ctx, givenNamespace)

			// when
			testEnvironment.EnsureK8sResourceCreated(t, ctx, tc.givenCompanion)

			// then
			_, err := testEnvironment.GetCompanionCRFromK8s(ctx, tc.givenCompanion.GetName(),
				tc.givenCompanion.GetNamespace())
			require.NoError(t, err)
		})
	}
}
