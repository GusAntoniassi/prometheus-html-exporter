package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExporterMetricsRegistry(t *testing.T) {
	registry, err := getExporterMetricsRegistry()
	ok(t, err)
	assert(t, registry != nil, "metric registry should not be null")
}

func TestProbeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/probe", nil)

	assert(t, err == nil, "http.NewRequest should not return an error, this is likely a problem in the test itself. error: %s", err)

	config := testExporterConfig
	server := getTestServer("<div id=\"foobar\">1</div>")

	for i := range config.Targets {
		config.Targets[i].Address = server.URL
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, config)
	})

	handler.ServeHTTP(rr, req)

	assert(t, rr.Code == http.StatusOK, "response should be of HTTP %d status, got %d", http.StatusOK, rr.Code)
	assert(t, rr.Body.String() != "", "response body should not be empty")
}
