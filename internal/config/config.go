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

	DSPTimeout             time.Duration
	DSPMaxIdleConns        int
	DSPMaxIdleConnsPerHost int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MaxBodyBytes int64
}

func Load() (Config, error) {
	var cfg Config
	var err error

	cfg.Addr = getString("ADDR", ":8080")
	cfg.PartnersPath = getString("PARTNERS_PATH", "partners.json")

	if cfg.AuctionTimeout, err = getDuration("AUCTION_TIMEOUT", 200*time.Millisecond); err != nil {
		return Config{}, err
	}
	if cfg.MaxConcurrency, err = getInt("MAX_CONCURRENCY", 50); err != nil {
		return Config{}, err
	}

	if cfg.DSPTimeout, err = getDuration("DSP_TIMEOUT", 3*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.DSPMaxIdleConns, err = getInt("DSP_MAX_IDLE_CONNS", 100); err != nil {
		return Config{}, err
	}
	if cfg.DSPMaxIdleConnsPerHost, err = getInt("DSP_MAX_IDLE_CONNS_PER_HOST", 20); err != nil {
		return Config{}, err
	}

	if cfg.ReadTimeout, err = getDuration("READ_TIMEOUT", 5*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.WriteTimeout, err = getDuration("WRITE_TIMEOUT", 10*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.IdleTimeout, err = getDuration("IDLE_TIMEOUT", 60*time.Second); err != nil {
		return Config{}, err
	}

	maxBody, err := getInt("MAX_BODY_BYTES", 1<<20)
	if err != nil {
		return Config{}, err
	}
	cfg.MaxBodyBytes = int64(maxBody)

	return cfg, nil
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
	if d <= 0 {
		return 0, fmt.Errorf("%s должен быть больше 0, получено %s", key, d)
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
	if n <= 0 {
		return 0, fmt.Errorf("%s должен быть больше 0, получено %d", key, n)
	}
	return n, nil
}
