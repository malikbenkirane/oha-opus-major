package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/malikbenkirane/oha-opus-major/service"
)

const defaultConfigDir = "/etc/player-data-service"

func main() {
	configDir := flag.String("config-dir", "", "path where to look for config dir")
	flag.Parse()

	if len(*configDir) == 0 {
		slog.Warn("Config directory flag not provided; using default", "configDir", defaultConfigDir)
		*configDir = defaultConfigDir
	}

	sv := service.NewMocker(*configDir)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := make(chan error, 1)
	defer close(err)

	go func() {
		err <- sv.Run(ctx)
	}()

	select {
	case <-quit:
		return
	case err := <-err:
		fmt.Println(err)
		os.Exit(1)
	}

}
