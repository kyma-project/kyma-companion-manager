package backendmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
	kcorev1 "k8s.io/api/core/v1"
)

func Test_getContainerPorts(t *testing.T) {
	t.Parallel()

	// when
	got := getContainerPorts()

	// then
	want := []kcorev1.ContainerPort{
		{
			Name:          backendPortName,
			ContainerPort: backendPortNum,
		},
		{
			Name:          backendMetricsPortName,
			ContainerPort: backendMetricsPortNum,
		},
	}
	require.Equal(t, want, got)
}
