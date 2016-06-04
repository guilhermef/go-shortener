package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	handler := &RedirectHandler{client: client}
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

}

func TestWhenURLDoesntExist(t *testing.T) {

	client := getRedisClient()

	handler := &RedirectHandler{client: client}
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

}
