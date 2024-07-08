package controller

// **NOTE:** This file contains unit tests for companion_controller.go.
// Integration tests for controller are located in the test/integration directory.

import (
	"testing"

	"github.com/stretchr/testify/require"

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
