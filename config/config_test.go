package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	settings := &readSettings{}
	options := getRedisOptions(settings)

	if options.Addr != "127.0.0.1:6379" {
		t.Fatalf("REDIS_HOST default value incorrect. Expected 127.0.0.1:6379; got %s", options.Addr)
	}
	if options.Password != "" {
		t.Fatalf("REDIS_PASS default value incorrect. Expected empty; got %s", options.Password)
	}
	if options.DB != 0 {
		t.Fatalf("REDIS_DB default value incorrect. Expected 0; got %d", options.DB)
	}

	out := getLoggerOutput(settings)

	if out != os.Stdout {
		t.Fatal("Default logger output should be STDOUT")
	}

	cfg, _ := NewConfig()

	if cfg.Port != "1234" {
		t.Fatalf("Unexpected port: %s", cfg.Port)
	}
}

func TestEnvVariables(t *testing.T) {
	os.Setenv("REDIS_HOST", "test-redisHost")
	os.Setenv("REDIS_PASS", "test-redisPass")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("PORT", "test-port")
	os.Setenv("LOG_PATH", "test-logPath")

	settings := &readSettings{
		LogPath:   "wrong-logPath",
		Port:      "wrong-port",
		RedisHost: "wrong-redisHost",
		RedisPass: "wrong-redisPass",
		RedisDB:   2}
	defer os.Remove("./wrong-logPath")

	options := getRedisOptions(settings)
	if options.Addr != "test-redisHost" {
		t.Fatalf("REDIS_HOST default value incorrect. Expected test-redisHost; got %s", options.Addr)
	}
	if options.Password != "test-redisPass" {
		t.Fatalf("REDIS_PASS default value incorrect. Expected test-redisPass; got %s", options.Password)
	}
	if options.DB != 1 {
		t.Fatalf("REDIS_DB default value incorrect. Expected 1; got %d", options.DB)
	}

	out := getLoggerOutput(settings)
	if out == os.Stdout {
		t.Fatal("Logger output should not be STDOUT")
	}

	cfg, _ := NewConfig()
	if cfg.Port != "test-port" {
		t.Fatalf("Unexpected port: %s", cfg.Port)
	}
}

func TestFetchFromValidYAML(t *testing.T) {
	settings := &readSettings{
		LogPath:   "test-logPath",
		Port:      "test-port",
		RedisHost: "test-redisHost",
		RedisPass: "test-redisPass",
		RedisDB:   1}
	b, _ := yaml.Marshal(settings)
	ioutil.WriteFile("./settings.yml", b, 0644)
	defer os.Remove("./settings.yml")
	defer os.Remove("./test-logPath")

	options := getRedisOptions(settings)
	if options.Addr != "test-redisHost" {
		t.Fatalf("REDIS_HOST default value incorrect. Expected test-redisHost; got %s", options.Addr)
	}
	if options.Password != "test-redisPass" {
		t.Fatalf("REDIS_PASS default value incorrect. Expected test-redisPass; got %s", options.Password)
	}
	if options.DB != 1 {
		t.Fatalf("REDIS_DB default value incorrect. Expected 1; got %d", options.DB)
	}

	out := getLoggerOutput(settings)
	if out == os.Stdout {
		t.Fatal("Logger output should not be STDOUT")
	}

	cfg, _ := NewConfig()
	if cfg.Port != "test-port" {
		t.Fatalf("Unexpected port: %s", cfg.Port)
	}
}

func TestFetchFromInvalidYAML(t *testing.T) {
	b := []byte("invalid\n\tYAML\nfile")
	ioutil.WriteFile("./settings.yml", b, 0644)
	defer os.Remove("./settings.yml")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("Invalid YAML did not cause error return")
	}
}
