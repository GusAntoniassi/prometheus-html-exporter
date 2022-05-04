package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestProbeHandlerWithQueryParams(t *testing.T) {
	config := getDefaultConfig()
	server := getTestServer("<div id=\"foobar\">1</div>")

	queryParams := url.Values{
		"selector":                []string{"//div[@id='foobar']/text()"},
		"decimal_point_separator": []string{","},
		"thousands_separator":     []string{"."},
		"metric_name":             []string{"wikipedia_articles_total"},
		"metric_type":             []string{"gauge"},
		"metric_help":             []string{"Total of articles available at Wikipedia"},
		"target":                  []string{server.URL},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/probe?%s", queryParams.Encode()), nil)

	assert(t, err == nil, "http.NewRequest should not return an error, this is likely a problem in the test itself. error: %s", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, config)
	})

	handler.ServeHTTP(rr, req)

	assert(t, rr.Code == http.StatusOK, "response should be of HTTP %d status, got %d", http.StatusOK, rr.Code)
	assert(t, rr.Body.String() != "", "response body should not be empty")
}

func TestProbeHandlerWithQueryParamsNotSupplied(t *testing.T) {
	// @TODO
}

func TestProbeHandlerWithQueryParamsMissingRequired(t *testing.T) {
	// @TODO
}
