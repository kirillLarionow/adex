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

	if auctionTimeout <= 0 {
		return Config{}, fmt.Errorf("AUCTION_TIMEOUT должен быть больше 0, получено %s", auctionTimeout)
	}

	dspTimeout, err := getDuration("DSP_TIMEOUT", 3*time.Second)
	if err != nil {
		return Config{}, err
	}

	if dspTimeout <= 0 {
		return Config{}, fmt.Errorf("DSP_TIMEOUT должен быть больше 0, получено %s", dspTimeout)
	}

	maxConcurrency, err := getInt("MAX_CONCURRENCY", 50)
	if err != nil {
		return Config{}, err
	}

	if maxConcurrency <= 0 {
		return Config{}, fmt.Errorf("MAX_CONCURRENCY должен быть больше 0, получено %d", maxConcurrency)
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
