package main

import (
	"testing"
)

func TestWithoutFlags(t *testing.T) {
	cfg := &config{}
	setConfig(cfg)
	if cfg.redisHost != "localhost" {
		t.Fatal("REDIS_HOST default value missing")
	}

}
