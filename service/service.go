package service

import (
	"context"
	"log/slog"

	"github.com/malikbenkirane/oha-opus-major/internal/adapter/player_data_repository/mock"
	"github.com/malikbenkirane/oha-opus-major/internal/adapter/player_data_server/http"
	"github.com/malikbenkirane/oha-opus-major/internal/port"
	"github.com/spf13/viper"
)

// NewMocker creates a mock Service using configuration from configDir.
// Falls back to the default address if not set in the config.
func NewMocker(configDir string) Service {
	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)

	addr := viper.GetString("server.addr")
	if addr == "" {
		slog.Warn("address missing in configuration, falling back to default", "addr", ":8080")
		addr = ":8080"
	}

	svc := &service{}
	svc.server = http.New(addr, mock.New())
	return svc
}

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	server port.PlayerDataServer
}

func (s service) Run(ctx context.Context) error {
	return s.server.Serve(ctx)
}
