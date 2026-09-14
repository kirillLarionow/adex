package auction

import (
	"adex/internal/domain"
	"adex/internal/pretargeting"
	"context"
	"log"
	"sync"
	"time"
)

type PartnerListner interface {
	List() []domain.Partner
}

type DSPClient interface {
	SendRequest(ctx context.Context, endpoint string, req domain.AuctionRequest) error
}

type Service struct {
	partners  PartnerListner
	dspClient DSPClient
}

func NewService(partners PartnerListner, dspClient DSPClient) *Service {
	return &Service{
		partners:  partners,
		dspClient: dspClient,
	}
}

func (s *Service) RunAuction(ctx context.Context, auctionRequest domain.AuctionRequest) domain.AuctionResponse {
	start := time.Now()
	partners := s.partners.List()
	matcheds := make([]domain.Partner, 0, len(partners))

	for _, partner := range partners {
		ok, reason := pretargeting.Match(auctionRequest, partner)

		if !ok {
			log.Printf("партнер: %s, причина не пройденной фильтрации: %s", partner.Name, reason)
			continue
		}

		matcheds = append(matcheds, partner)
	}

	succeeded := s.sendToPartners(ctx, auctionRequest, matcheds)

	names := make([]string, len(matcheds))

	for i, name := range matcheds {
		names[i] = name.Name
	}

	return domain.AuctionResponse{
		RequestID:   auctionRequest.RequestID,
		MatchedDSPs: names,
		Sent:        len(names),
		Succeeded:   succeeded,
		DurationMs:  time.Since(start).Milliseconds(),
	}
}

func (s *Service) sendToPartners(
	ctx context.Context,
	auctionRequest domain.AuctionRequest,
	matcheds []domain.Partner) int {

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)

	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	succeeded := 0

	for _, domainPartner := range matcheds {

		wg.Add(1)

		go func(domainPartner domain.Partner) {
			defer wg.Done()

			if err := s.dspClient.SendRequest(
				ctx,
				domainPartner.EndPoint,
				auctionRequest); err != nil {
				log.Printf("партнер: %s, упал: %v", domainPartner.Name, err)
				return
			}

			mu.Lock()
			succeeded++
			mu.Unlock()

		}(domainPartner)
	}

	wg.Wait()

	return succeeded
}
