package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"strings"
	"gopkg.in/redis.v3"
)

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func TestWhenURLExists(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()

	client.Del("go-shortener-count:/location-test")

	err := client.Set(
		"go-shortener:/location-test",
		"http://location.test",
		0,
	).Err()

	if err != nil {
		t.Fatal(err)
	}

	handler := &RedirectHandler{client: client, logger: logger}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/location-test", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)

	count, _ := client.Get("go-shortener-count:/location-test").Result()

	if count != "1" {
		t.Fatal("Unexpected count")
	}

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 301 {
		t.Fatalf("Received non-301 response: %d\n", resp.StatusCode)
	}
	if resp.Header.Get("Location") != "http://location.test" {
		t.Errorf("Received wrong location: %s\n", resp.Header.Get("Location"))
	}

	if !strings.HasSuffix(buffer.String(), "301 /location-test http://location.test\n") {
		t.Fatalf("Incorrect log entry: %s\n", buffer.String())
	}
}

func TestWhenURLDoesntExist(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()

	handler := &RedirectHandler{client: client, logger: logger}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/should-not-exist", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Received non-404 response: %d\n", resp.StatusCode)
	}

	if !strings.HasSuffix(buffer.String(), "404 /should-not-exist\n") {
		t.Fatalf("Incorrect log entry: %s\n", buffer.String())
	}

}
