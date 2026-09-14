package pretargeting

import (
	"adex/internal/domain"
	"slices"
)

func Match(auctionRequest domain.AuctionRequest, partner domain.Partner) (bool, string) {
	if !partner.IsEnabled {
		return false, "партнер дизейбл"
	}

	if len(partner.Countries) > 0 && !slices.Contains(partner.Countries, auctionRequest.Country) {
		return false, "страны не подходят"
	}

	if len(partner.DeviceTypes) > 0 && !slices.Contains(partner.DeviceTypes, auctionRequest.DeviceType) {
		return false, "девайсы не подходят"
	}

	if partner.MinBidFloor > auctionRequest.BidFloor {
		return false, "ниже минимальной ставки"
	}

	for _, category := range auctionRequest.Categories {
		if slices.Contains(partner.BlockedCategories, category) {
			return false, "категория заблочена " + category
		}
	}

	return true, ""
}
