package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr           string
	PartnersPath   string
	AuctionTimeout time.Duration
	MaxConcurrency int
	DSPTimeout     time.Duration
}

func Load() (Config, error) {
	auctionTimeout, err := getDuration("AUCTION_TIMEOUT", 200*time.Millisecond)
	if err != nil {
		return Config{}, err
	}

	dspTimeout, err := getDuration("DSP_TIMEOUT", 3*time.Second)
	if err != nil {
		return Config{}, err
	}

	maxConcurrency, err := getInt("MAX_CONCURRENCY", 50)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Addr:           getString("ADDR", ":8080"),
		PartnersPath:   getString("PARTNERS_PATH", "partners.json"),
		AuctionTimeout: auctionTimeout,
		MaxConcurrency: maxConcurrency,
		DSPTimeout:     dspTimeout,
	}, nil
}

func getString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return d, nil
}

func getInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return n, nil
}
