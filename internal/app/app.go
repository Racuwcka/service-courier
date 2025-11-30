package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/dig"

	"github.com/Racuwcka/service-courier/internal/config"
	"github.com/Racuwcka/service-courier/internal/http-server/handlers/courier/add"
	"github.com/Racuwcka/service-courier/internal/http-server/handlers/courier/provider"
	"github.com/Racuwcka/service-courier/internal/http-server/handlers/courier/router"
	"github.com/Racuwcka/service-courier/internal/kafka"
)

func NewApp() {
	c := dig.New()

	_ = c.Provide(provideConfig)
	// todo: logger
	_ = c.Provide(provideConsumer)
	// todo: как-то переделать
	_ = c.Provide(func() add.Provider {
		return provider.NewCourierProvider()
	})
	_ = c.Provide(add.New)
	_ = c.Provide(router.New)
	_ = c.Provide(provideRouter)

	err := c.Invoke(func(r *chi.Mux, cfg *config.Config, h *router.CourierHandlers, consumer *kafka.Consumer) error {
		r.Mount("/courier", h.Routes())

		srv := &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.App.Port),
			Handler: r,
		}

		shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// ---- START KAFKA ----
		go consumer.Start(shutdownCtx)

		// ---- START HTTP ----
		go func() {
			log.Println("service started on port", cfg.App.Port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen and serve: %v", err)
			}
		}()

		// ---- WAIT ----
		<-shutdownCtx.Done()
		log.Println("shutting down gracefully...")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = consumer.Close()
		if err := srv.Shutdown(ctxTimeout); err != nil {
			log.Fatalf("server forced to shutdown: %v", err)
		}

		log.Println("server stopped")
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func provideConfig() *config.Config {
	return config.MustLoad()
}

func provideRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return r
}

func provideConsumer(cfg *config.Config) *kafka.Consumer {
	c, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.OrderTopic,
		"courier-service-group",
	)

	if err == nil {
		return c
	}

	log.Printf("failed to init sarama producer, fallback to noop: %v", err)

	return &kafka.Consumer{}
}
