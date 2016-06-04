package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyHandler(t *testing.T) {
	handler := &RedirectHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
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
