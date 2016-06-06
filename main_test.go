package main

import (
	"testing"
)

func TestWithoutFlags(t *testing.T) {
	cfg := &config{}
	setConfig(cfg)
	if cfg.redisHost != "localhost:6379" {
		t.Fatal("REDIS_HOST default value missing")
	}

	if cfg.logPath != "" {
		t.Fatalf("Unexpected logPath %s", cfg.logPath)
	}
}
