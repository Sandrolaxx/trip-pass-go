package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("goodbye :)")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	
	logger, err := cfg.Build()

	if err != nil {
		return err
	}

	logger = logger.Named("Trip_pass_app")

	defer func ()  { _ = logger.Sync()} ()

	// pool, err := pgxpool.New(ctx, fmt.Sprintf())

	select {
	case <- ctx.Done():
	return nil
	}
}