package pretargeting

import (
	"adex/internal/domain"
	"testing"
)

func TestMatch(t *testing.T) {
	partner := domain.Partner{
		Name:              "DSP Alpha",
		IsEnabled:         true,
		Countries:         []string{"RU", "KZ"},
		DeviceTypes:       []string{"mobile"},
		MinBidFloor:       1.0,
		BlockedCategories: []string{"gambling"},
	}

	tests := []struct {
		name       string
		req        domain.AuctionRequest
		partner    domain.Partner
		wantOk     bool
		wantReason string
	}{
		{
			name: "всё подходит",
			req: domain.AuctionRequest{
				Country: "RU", DeviceType: "mobile", BidFloor: 2.0, Categories: []string{"news"},
			},
			partner: partner,
			wantOk:  true,
		},
		{
			name: "партнёр выключен",
			req: domain.AuctionRequest{
				Country: "RU", DeviceType: "mobile", BidFloor: 2.0,
			},
			partner:    domain.Partner{IsEnabled: false},
			wantOk:     false,
			wantReason: "партнер дизейбл",
		},
		{
			name: "страна не подходит",
			req: domain.AuctionRequest{
				Country: "DE", DeviceType: "mobile", BidFloor: 2.0,
			},
			partner:    partner,
			wantOk:     false,
			wantReason: "страны не подходят",
		},
		{
			name: "устройство не подходит",
			req: domain.AuctionRequest{
				Country: "RU", DeviceType: "tv", BidFloor: 2.0,
			},
			partner:    partner,
			wantOk:     false,
			wantReason: "девайсы не подходят",
		},
		{
			name: "ставка ниже минимума",
			req: domain.AuctionRequest{
				Country: "RU", DeviceType: "mobile", BidFloor: 0.5,
			},
			partner:    partner,
			wantOk:     false,
			wantReason: "ниже минимальной ставки",
		},
		{
			name: "заблокированная категория",
			req: domain.AuctionRequest{
				Country: "RU", DeviceType: "mobile", BidFloor: 2.0, Categories: []string{"gambling"},
			},
			partner:    partner,
			wantOk:     false,
			wantReason: "категория заблочена gambling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, reason := Match(tt.req, tt.partner)
			if ok != tt.wantOk {
				t.Errorf("функция вернула = %v, что мы ожидаем %v", ok, tt.wantOk)
			}
			if !tt.wantOk && reason != tt.wantReason {
				t.Errorf("причина = %q, ожидаемая причина %q", reason, tt.wantReason)
			}
		})
	}
}
