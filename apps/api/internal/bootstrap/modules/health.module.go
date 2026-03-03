package modules

import "authentication-project-exam/internal/adapter/inbound/web/handler"

type HealthModule struct {
	Handler *handler.HealthHandler
}

func NewHealthModule() *HealthModule {
	return &HealthModule{
		Handler: handler.NewHealthHandler(),
	}
}
