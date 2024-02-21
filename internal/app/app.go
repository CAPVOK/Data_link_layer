package app

import (
	"context"

	"github.com/lud0m4n/Network/internal/config"
	"github.com/lud0m4n/Network/internal/http/api"
)

// Application представляет основное приложение.
type Application struct {
	Config  *config.Config
	Handler *api.Handler
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
	// Инициализируйте конфигурацию
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	h := api.NewHandler()
	// Инициализируйте и настройте объект Application
	app := &Application{
		Config:  cfg,
		Handler: h,
	}

	return app, nil
}
