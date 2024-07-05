package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"github.com/kyma-project/kyma-companion-manager/test/utils"
)

func Test_containsFinalizer(t *testing.T) {
	t.Parallel()

	// define test cases
	testCases := []struct {
		name           string
		givenCompanion *kcmv1alpha1.Companion
		wantResult     bool
	}{
		{
			name:           "should return false when finalizer is missing",
			givenCompanion: utils.NewCompanionCR(),
			wantResult:     false,
		},
		{
			name:           "should return true when finalizer is present",
			givenCompanion: utils.NewCompanionCR(utils.WithCompanionCRFinalizer(FinalizerName)),
			wantResult:     true,
		},
	}

	// run test cases
	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			// given
			testEnv := NewMockedUnitTestEnvironment(t)
			reconciler := testEnv.Reconciler

			// when, then
			require.Equal(t, testcase.wantResult, reconciler.containsFinalizer(testcase.givenCompanion))
		})
	}
}

func Test_addFinalizer(t *testing.T) {
	// given
	givenCompanion := utils.NewCompanionCR()

	testEnv := NewMockedUnitTestEnvironment(t, givenCompanion)
	reconciler := testEnv.Reconciler

	// when
	_, err := reconciler.addFinalizer(context.Background(), givenCompanion)

	// then
	require.NoError(t, err)
	gotCompanion, err := testEnv.GetCompanion(givenCompanion.GetName(), givenCompanion.GetNamespace())
	require.NoError(t, err)
	require.True(t, reconciler.containsFinalizer(&gotCompanion))
}

func Test_removeFinalizer(t *testing.T) {
	// given
	givenCompanion := utils.NewCompanionCR(utils.WithCompanionCRFinalizer(FinalizerName))

	testEnv := NewMockedUnitTestEnvironment(t, givenCompanion)
	reconciler := testEnv.Reconciler

	// when
	_, err := reconciler.removeFinalizer(context.Background(), givenCompanion)

	// then
	require.NoError(t, err)
	gotCompanion, err := testEnv.GetCompanion(givenCompanion.GetName(), givenCompanion.GetNamespace())
	require.NoError(t, err)
	require.False(t, reconciler.containsFinalizer(&gotCompanion))
}
