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

	if options.Addr != "localhost:6379" {
		t.Fatalf("REDIS_HOST default value incorrect. Expected localhost:6379; got %s", options.Addr)
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
  os.Setenv("PORT", "test-port")
  os.Setenv("LOG_PATH", "test-logPath")

  settings := &readSettings{LogPath: "wrong-logPath", Port: "wrong-port", RedisHost: "wrong-redisHost"}
  defer os.Remove("./wrong-logPath")

  options := getRedisOptions(settings)
  if options.Addr != "test-redisHost" {
    t.Fatalf("REDIS_HOST default value incorrect. Expected test-redisHost; got %s", options.Addr)
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
	settings := &readSettings{LogPath: "test-logPath", Port: "test-port", RedisHost: "test-redisHost"}
	b, _ := yaml.Marshal(settings)
	ioutil.WriteFile("./settings.yml", b, 0644)
	defer os.Remove("./settings.yml")
	defer os.Remove("./test-logPath")

	options := getRedisOptions(settings)
	if options.Addr != "test-redisHost" {
		t.Fatalf("REDIS_HOST default value incorrect. Expected test-redisHost; got %s", options.Addr)
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
