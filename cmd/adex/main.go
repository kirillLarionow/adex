package main

import (
	"adex/internal/auction"
	"adex/internal/config"
	"adex/internal/dsp"
	"adex/internal/httpapi"
	"adex/internal/partner"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	repo, err := partner.NewRepository(cfg.PartnersPath)
	if err != nil {
		log.Fatal(err)
	}

	client := dsp.NewHttpClient(cfg.DSPTimeout, cfg.DSPMaxIdleConns, cfg.DSPMaxIdleConnsPerHost)

	service := auction.NewService(repo, client, cfg.AuctionTimeout, cfg.MaxConcurrency)

	handler := httpapi.NewHandler(service, cfg.MaxBodyBytes)

	mux := http.NewServeMux()
	mux.HandleFunc("/auction", handler.HandleAuction)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("adex: listening on", cfg.Addr)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("завершаю работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("сервер остановлен")
}
