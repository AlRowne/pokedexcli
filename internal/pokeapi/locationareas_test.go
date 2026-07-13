package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationAreas_ParsesResponse(t *testing.T) {
	previous := "https://example.com/prev"
	fixture := LocationAreaResponse{
		Count:    2,
		Next:     nil,
		Previous: &previous,
		Results: []LocationArea{
			{Name: "canalave-city-area", URL: "https://example.com/1"},
			{Name: "eterna-city-area", URL: "https://example.com/2"},
		},
	}
	data, err := json.Marshal(fixture)
	if err != nil {
		t.Fatalf("failed to marshal test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer server.Close()

	got, err := GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Count != fixture.Count {
		t.Errorf("expected count %d, got %d", fixture.Count, got.Count)
	}
	if got.Next != nil {
		t.Errorf("expected Next to be nil, got %q", *got.Next)
	}
	if got.Previous == nil || *got.Previous != previous {
		t.Errorf("expected Previous %q, got %v", previous, got.Previous)
	}
	if len(got.Results) != len(fixture.Results) {
		t.Fatalf("expected %d results, got %d", len(fixture.Results), len(got.Results))
	}
	for i, area := range got.Results {
		if area != fixture.Results[i] {
			t.Errorf("result %d: expected %+v, got %+v", i, fixture.Results[i], area)
		}
	}
}

func TestGetLocationAreas_UsesCache(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Write([]byte(`{"count":0,"next":null,"previous":null,"results":[]}`))
	}))
	defer server.Close()

	if _, err := GetLocationAreas(server.URL); err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	if _, err := GetLocationAreas(server.URL); err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request to reach the server (second call should be served from cache), got %d", requestCount)
	}
}

func TestGetLocationAreas_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	if _, err := GetLocationAreas(server.URL); err == nil {
		t.Error("expected an error for invalid JSON response, got nil")
	}
}
