package main

import (
	"github.com/paulrozhkin/metrics-server/config"
	"github.com/paulrozhkin/metrics-server/internal/http_server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	fx.New(
		fx.Provide(
			config.InitLogger,
			zap.L,
			zap.S,
			config.LoadConfigurations,
			http_server.NewServerRoute,
			http_server.NewHTTPServer,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(func(*config.LoggerConfigurator) {}),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
