package handler

import (
	"net/http"
	"gopkg.in/redis.v3"
	"net/http/httptest"
	"testing"
)

func TestMyHandler(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Set(
		"go-shortner:/location-test",
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
