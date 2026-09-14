package domain

type AuctionRequest struct {
	RequestID  string   `json:"request_id"`
	Country    string   `json:"country"`
	DeviceType string   `json:"device_type"`
	BidFloor   float64  `json:"bid_floor"`
	Categories []string `json:"categories"`
}

type AuctionResponse struct {
	RequestID   string   `json:"request_id"`
	MatchedDSPs []string `json:"matched_dsps"`
	Sent        int      `json:"sent"`
	Succeeded   int      `json:"succeeded"`
	DurationMs  int64    `json:"duration_ms"`
}

type Partner struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	EndPoint          string   `json:"endpoint"`
	IsEnabled         bool     `json:"is_enabled"`
	Countries         []string `json:"countries"`
	DeviceTypes       []string `json:"device_types"`
	MinBidFloor       float64  `json:"min_bid_floor"`
	BlockedCategories []string `json:"blocked_categories"`
}
