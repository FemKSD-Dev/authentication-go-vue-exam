package bootstrap

import "authentication-project-exam/internal/adapter/inbound/web"

func NewHTTPServer(container *Container) *web.ServerDependencies {
	if container == nil {
		panic("container is nil")
	}

	return &web.ServerDependencies{
		AuthHandler:         container.AuthHandler,
		AuthMiddleware:      container.AuthMiddleware,
		HealthHandler:       container.HealthHandler,
		RequestIDMiddleware: container.RequestIDMiddleware,
		LogMiddleware:       container.LogMiddleware,
	}
}
