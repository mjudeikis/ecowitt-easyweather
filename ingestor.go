package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mjudeikis/weewx-easyweather/pkg/config"
	"github.com/mjudeikis/weewx-easyweather/pkg/server"
	"github.com/mjudeikis/weewx-easyweather/pkg/utils/ratelimiter/logger"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Printf("error starting controller: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	c, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.GetLoggerInstance("", logger.ParseLogLevel(c.LogLevel))
	log.Info("starting ingestor")

	server, err := server.New(log, c)
	if err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go server.Run(ctx)
	select {
	case <-signals:
		// shutdown
	case <-ctx.Done():
		// ctx termination
	}

	return nil
}
