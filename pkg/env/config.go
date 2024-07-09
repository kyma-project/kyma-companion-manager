package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the environment config for kyma companion manager.
type Config struct {
	// KymaCompanionBackendImage container image for kyma-companion-backend.
	KymaCompanionBackendImage string `envconfig:"KYMA_COMPANION_BACKEND_IMAGE" required:"true"`
}

func GetConfig() Config {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	return cfg
}
