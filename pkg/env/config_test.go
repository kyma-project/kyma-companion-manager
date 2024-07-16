package env

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GetConfig(t *testing.T) {
	g := NewGomegaWithT(t)
	envs := map[string]string{
		// required
		"KYMA_COMPANION_BACKEND_IMAGE": "test:latest",
	}

	for k, v := range envs {
		t.Setenv(k, v)
	}
	config := GetConfig()
	// Ensure required variables can be set
	g.Expect(config.KymaCompanionBackendImage).To(Equal(envs["KYMA_COMPANION_BACKEND_IMAGE"]))
}
