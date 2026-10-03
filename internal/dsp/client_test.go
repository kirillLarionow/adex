package dsp

import (
	"adex/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendRequest_OK(t *testing.T) {
	var got domain.AuctionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("метод = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type = %q, want application/json", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("не удалось разобрать тело: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHttpClient(time.Second, 10, 10)
	req := domain.AuctionRequest{RequestID: "req-1", Country: "RU"}

	if err := client.SendRequest(context.Background(), server.URL, req); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if got.RequestID != "req-1" {
		t.Errorf("request_id = %q, want req-1", got.RequestID)
	}
}

func TestSendRequest_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHttpClient(time.Second, 10, 10)

	if err := client.SendRequest(context.Background(), server.URL, domain.AuctionRequest{}); err == nil {
		t.Fatal("ожидалась ошибка на статус 500, получили nil")
	}
}

func TestSendRequest_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(500 * time.Millisecond):
		case <-r.Context().Done():
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHttpClient(time.Second, 10, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := client.SendRequest(ctx, server.URL, domain.AuctionRequest{})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("запрос длился %v, клиент не прервал его по контексту", elapsed)
	}
}
