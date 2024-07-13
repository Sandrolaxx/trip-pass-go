package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trip-pass-go/internal/api"
	"trip-pass-go/internal/api/spec"
	"trip-pass-go/internal/mailer/mailpit"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//Criando contexto da aplicação
	ctx := context.Background()
	//Caso ocorra algum desses comando vai cancelar o contexto encerrando a aplicação
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("goodbye :)")
}

func run(ctx context.Context) error {
	//configurando logger
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()

	if err != nil {
		return err
	}

	logger = logger.Named("Trip_pass_app")

	defer func() { _ = logger.Sync() }()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("JOURNEY_DATABASE_USER"),
		os.Getenv("JOURNEY_DATABASE_PASSWORD"),
		os.Getenv("JOURNEY_DATABASE_HOST"),
		os.Getenv("JOURNEY_DATABASE_PORT"),
		os.Getenv("JOURNEY_DATABASE_NAME"),
	),
	)

	if err != nil {
		return err
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	//Definindo interface do server
	serverInterface := api.NewAPI(pool, logger, mailpit.NewMailTrip(pool))

	//Criando roteamento
	router := chi.NewMux()
	router.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))

	router.Mount("/", spec.Handler(serverInterface))

	//Definindo dados do server
	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	//Gracefull shutdown para encerrar conexões do banco
	defer func() {
		const timeout = 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown server", zap.Error(err))
		}

	}()

	//Criando canal de erro, para termos uma comunicação assincrona no sistema
	errorChannel := make(chan error, 1)

	//Criando meu servidor em uma GoRoutine pois o método ListenAndServe é blocante,
	//se chamar na rotina principal travaria o serviço
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errorChannel <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errorChannel:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
