package main

import (
	"adex/internal/auction"
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
	repo, err := partner.NewRepository("partners.json")
	if err != nil {
		log.Fatal(err)
	}

	client := dsp.NewHttpClient()

	service := auction.NewService(repo, client)

	handler := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/auction", handler.HandleAuction)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("adex: listening on :8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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
