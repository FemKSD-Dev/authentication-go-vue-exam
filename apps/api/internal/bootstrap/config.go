package bootstrap

import "authentication-project-exam/internal/config"

func LoadConfig() (*config.Config, error) {
	return config.Load()
}
