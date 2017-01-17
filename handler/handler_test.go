package handler

import (
	"bytes"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	handler := &RedirectHandler{Client: client, Logger: logger}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/location-test", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()

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

func TestWhenURLDoesntExistAndDoesntExistRedirectHost(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()

	handler := &RedirectHandler{Client: client, Logger: logger}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/should-not-exist", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()

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

func TestWhenURLDoesntExistAndExistRedirectHost(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()
	extra := getExtra(&Extra{RedirectHost: "http://luke.wars"})

	handler := &RedirectHandler{Client: client, Logger: logger, Extra: extra}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/should-not-exist", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 302 {
		t.Errorf("Received non-302 response: %d\n", resp.StatusCode)
	}

	if resp.Header.Get("Location") != "http://luke.wars" {
		t.Errorf("Received wrong location: %s\n", resp.Header.Get("Location"))
	}

	if !strings.HasSuffix(buffer.String(), "302 /should-not-exist http://luke.wars\n") {
		t.Fatalf("Incorrect log entry: %s\n", buffer.String())
	}
}

func TestWhenURLDoesntExistAndExistRedirectHostAndExistRedirectCode(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()
	extra := getExtra(&Extra{RedirectHost: "http://luke.wars", RedirectCode: 200})

	handler := &RedirectHandler{Client: client, Logger: logger, Extra: extra}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/should-not-exist", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Received non-200 response: %d\n", resp.StatusCode)
	}

	if resp.Header.Get("Location") != "http://luke.wars" {
		t.Errorf("Received wrong location: %s\n", resp.Header.Get("Location"))
	}

	if !strings.HasSuffix(buffer.String(), "200 /should-not-exist http://luke.wars\n") {
		t.Fatalf("Incorrect log entry: %s\n", buffer.String())
	}
}

func TestHealthCheck(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	client := getRedisClient()

	handler := &RedirectHandler{Client: client, Logger: logger}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/healthcheck", nil)
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Received non-200 response: %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != "WORKING\n" {
		t.Errorf("Incorrect response body: %s\n", string(body))
	}
	if !strings.HasSuffix(buffer.String(), "200 /healthcheck\n") {
		t.Fatalf("Incorrect log entry: %s\n", buffer.String())
	}
}
