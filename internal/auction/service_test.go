package auction

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"adex/internal/domain"
)

type fakePartners struct {
	list []domain.Partner
}

func (f fakePartners) List() []domain.Partner {
	return f.list
}

type fakeDSP struct {
	mu          sync.Mutex
	calls       []string
	fail        map[string]bool
	delay       map[string]time.Duration
	inFlight    int
	maxInFlight int
}

func (f *fakeDSP) SendRequest(ctx context.Context, endpoint string, req domain.AuctionRequest) error {
	f.mu.Lock()
	f.calls = append(f.calls, endpoint)
	f.inFlight++
	if f.inFlight > f.maxInFlight {
		f.maxInFlight = f.inFlight
	}
	f.mu.Unlock()

	defer func() {
		f.mu.Lock()
		f.inFlight--
		f.mu.Unlock()
	}()

	if d := f.delay[endpoint]; d > 0 {
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if f.fail[endpoint] {
		return errors.New("dsp error")
	}
	return nil
}

func newPartner(name, endpoint string, countries ...string) domain.Partner {
	return domain.Partner{
		Name:      name,
		EndPoint:  endpoint,
		IsEnabled: true,
		Countries: countries,
	}
}

var baseRequest = domain.AuctionRequest{
	RequestID:  "req-1",
	Country:    "RU",
	DeviceType: "mobile",
	BidFloor:   1.0,
}

func TestRunAuction_SendsOnlyToMatchedPartners(t *testing.T) {
	disabled := newPartner("C", "c", "RU")
	disabled.IsEnabled = false

	partners := fakePartners{list: []domain.Partner{
		newPartner("A", "a", "RU"),
		newPartner("B", "b", "US"),
		disabled,
	}}
	dsp := &fakeDSP{}

	service := NewService(partners, dsp, time.Second, 10)
	resp := service.RunAuction(context.Background(), baseRequest)

	if len(dsp.calls) != 1 || dsp.calls[0] != "a" {
		t.Fatalf("запросы ушли на %v, want [a]", dsp.calls)
	}
	if len(resp.MatchedDSPs) != 1 || resp.MatchedDSPs[0] != "A" {
		t.Errorf("matched_dsps = %v, want [A]", resp.MatchedDSPs)
	}
	if resp.Sent != 1 {
		t.Errorf("sent = %d, want 1", resp.Sent)
	}
}

func TestRunAuction_CountsSucceeded(t *testing.T) {
	partners := fakePartners{list: []domain.Partner{
		newPartner("A", "a", "RU"),
		newPartner("B", "b", "RU"),
		newPartner("C", "c", "RU"),
	}}
	dsp := &fakeDSP{fail: map[string]bool{"b": true}}

	service := NewService(partners, dsp, time.Second, 10)
	resp := service.RunAuction(context.Background(), baseRequest)

	if resp.Sent != 3 {
		t.Errorf("sent = %d, want 3", resp.Sent)
	}
	if resp.Succeeded != 2 {
		t.Errorf("succeeded = %d, want 2", resp.Succeeded)
	}
}

func TestRunAuction_SlowPartnerTimesOut(t *testing.T) {
	partners := fakePartners{list: []domain.Partner{
		newPartner("A", "a", "RU"),
		newPartner("B", "b", "RU"),
	}}
	dsp := &fakeDSP{delay: map[string]time.Duration{"b": 500 * time.Millisecond}}

	service := NewService(partners, dsp, 50*time.Millisecond, 10)

	start := time.Now()
	resp := service.RunAuction(context.Background(), baseRequest)
	elapsed := time.Since(start)

	if resp.Sent != 2 {
		t.Errorf("sent = %d, want 2", resp.Sent)
	}
	if resp.Succeeded != 1 {
		t.Errorf("succeeded = %d, want 1", resp.Succeeded)
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("аукцион длился %v, не уложился в таймаут", elapsed)
	}
}

func TestRunAuction_LimitsConcurrency(t *testing.T) {
	var list []domain.Partner
	delay := map[string]time.Duration{}
	for i := 0; i < 10; i++ {
		endpoint := fmt.Sprintf("p%d", i)
		list = append(list, newPartner(endpoint, endpoint, "RU"))
		delay[endpoint] = 20 * time.Millisecond
	}
	dsp := &fakeDSP{delay: delay}

	service := NewService(fakePartners{list: list}, dsp, time.Second, 3)
	resp := service.RunAuction(context.Background(), baseRequest)

	if resp.Succeeded != 10 {
		t.Errorf("succeeded = %d, want 10", resp.Succeeded)
	}
	if dsp.maxInFlight > 3 {
		t.Errorf("одновременно было %d запросов, лимит 3", dsp.maxInFlight)
	}
	if dsp.maxInFlight < 2 {
		t.Errorf("одновременно было %d запросов, запросы шли не параллельно", dsp.maxInFlight)
	}
}
