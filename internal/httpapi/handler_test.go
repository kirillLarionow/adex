package httpapi

import (
	"adex/internal/domain"

	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRunner struct {
	called bool
}

func (f *fakeRunner) RunAuction(ctx context.Context, req domain.AuctionRequest) domain.AuctionResponse {
	f.called = true
	return domain.AuctionResponse{RequestID: req.RequestID}
}

const validBody = `{
"request_id":"req-1",
"country":"RU",
"device_type":"mobile",
"bid_floor":1.5,
"categories":["news"]
}`

func TestHandleAuction(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		wantStatus  int
		wantCalled  bool
	}{
		{
			name:        "валидный запрос",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        validBody,
			wantStatus:  http.StatusOK,
			wantCalled:  true,
		},
		{
			name:        "content-type с charset",
			method:      http.MethodPost,
			contentType: "application/json; charset=utf-8",
			body:        validBody,
			wantStatus:  http.StatusOK,
			wantCalled:  true,
		},
		{
			name:        "метод GET",
			method:      http.MethodGet,
			contentType: "application/json",
			body:        validBody,
			wantStatus:  http.StatusMethodNotAllowed,
		},
		{
			name:        "не json content-type",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        validBody,
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:       "без content-type",
			method:     http.MethodPost,
			body:       validBody,
			wantStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:        "битый json",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"request_id":`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "неизвестное поле",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"request_id":"req-1","country":"RU","device_type":"mobile","bidfloor":1.5}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "нет request_id",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"country":"RU","device_type":"mobile"}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "неизвестный device_type",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"request_id":"req-1","country":"RU","device_type":"fridge"}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "слишком большое тело",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"request_id":"` + strings.Repeat("a", 1<<20) + `"}`,
			wantStatus:  http.StatusRequestEntityTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeRunner{}
			handler := NewHandler(runner)
			req := httptest.NewRequest(tt.method, "/auction", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rec := httptest.NewRecorder()

			handler.HandleAuction(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("статус = %d, want %d, body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if runner.called != tt.wantCalled {
				t.Errorf("сервис вызван = %v, want %v", runner.called, tt.wantCalled)
			}
		})
	}
}
