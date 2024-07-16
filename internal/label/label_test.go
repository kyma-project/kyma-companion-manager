package label

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetCommonLabels(t *testing.T) {
	t.Parallel()

	// given
	name := "test-companion"
	// when
	got := GetCommonLabels(name)

	// then
	want := map[string]string{
		"app.kubernetes.io/component":  "companion",
		"app.kubernetes.io/created-by": "kyma-companion-manager",
		"app.kubernetes.io/instance":   "test-companion",
		"app.kubernetes.io/managed-by": "kyma-companion-manager",
		"app.kubernetes.io/name":       "test-companion",
		"app.kubernetes.io/part-of":    "test-companion",
		"kyma-project.io/dashboard":    "companion",
	}
	require.Equal(t, want, got)
}
