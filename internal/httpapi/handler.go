package httpapi

import (
	"adex/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuctionRunner interface {
	RunAuction(ctx context.Context, auctionRequest domain.AuctionRequest) domain.AuctionResponse
}

type Handler struct {
	service AuctionRunner
}

func NewHandler(service AuctionRunner) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleAuction(w http.ResponseWriter, r *http.Request) {

	var auctionRequest domain.AuctionRequest

	if r.Method != http.MethodPost {
		http.Error(w, "метод не подходит", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content-type должен быть application/json", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&auctionRequest); err != nil {
		http.Error(w, "json невалиден", http.StatusBadRequest)
		return
	}

	if err := validate(auctionRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.service.RunAuction(r.Context(), auctionRequest)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validate(req domain.AuctionRequest) error {
	if req.RequestID == "" {
		return fmt.Errorf("нужен RequestID")
	}

	if req.Country == "" {
		return fmt.Errorf("страна обязательна")
	}

	switch req.DeviceType {
	case "mobile", "desktop", "tv":
	default:
		return fmt.Errorf("девайс должен подходить")
	}

	if req.BidFloor < 0 {
		return fmt.Errorf("bid_floor должен быть >= 0")
	}

	return nil
}
