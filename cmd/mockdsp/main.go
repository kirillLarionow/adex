package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	startServer(9001, 0)
	startServer(9002, 0)
	startServer(9003, 500*time.Millisecond)

	select {}
}

func startServer(port int, delay time.Duration) {
	mux := http.NewServeMux()
	mux.HandleFunc("/bid", func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.WriteHeader(http.StatusOK)
	})
	go func() {
		addr := fmt.Sprintf(":%d", port)
		fmt.Println("запускаю сервер на", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatal(err)
		}
	}()
}
