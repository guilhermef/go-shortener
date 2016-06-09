package main

import (
	"os"
	"testing"
)

func TestWithoutFlags(t *testing.T) {
	cfg := &config{}
	setConfig(cfg)
	if cfg.redisHost != "localhost:6379" {
		t.Fatal("REDIS_HOST default value missing")
	}

	if cfg.logPath != "" {
		t.Fatalf("Unexpected logPath: %s", cfg.logPath)
	}

	if cfg.port != "1234" {
		t.Fatalf("Unexpected port: %s", cfg.port)
	}
}

func TestIfLoggerWritesToPath(t *testing.T) {
	cfg := &config{logPath: "./temp.log"}
	logger := getLogger(cfg)
	logger.Print("Test")

	err := os.Remove("./temp.log")
	if err != nil {
		t.Fatalf("Logger failed to write to given path: %s", err)
	}
}
