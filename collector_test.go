package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestCollect(t *testing.T) {
	html := "<div id=\"foobar\">1,234,567.08</div>"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, html)
	}))

	config := testExporterConfig
	config.ScrapeConfig.Address = server.URL

	collector := collector{config: config}

	ch := make(chan prometheus.Metric, 1)

	collector.Collect(ch)
	metric := <-ch
	assert(t, metric != nil, "collect should not return a nil metric")
}

func TestCollect_scrapeError(t *testing.T) {
	config := testExporterConfig
	config.ScrapeConfig.Address = "foo://bar.dev"

	collector := collector{config: config}

	defer func() { recover() }()

	ch := make(chan prometheus.Metric, 1)

	collector.Collect(ch)
	t.Errorf("should panic when configured with a summary metric (not implemented yet)")
}

func TestMakeNewConstMetric(t *testing.T) {
	value := 123.0
	_, err := makeNewConstMetric(testExporterConfig, value)

	ok(t, err)
}

func TestMakeMetricDesc(t *testing.T) {
	config := testExporterConfig
	expected := config.GlobalConfig.MetricNamePrefix + config.ScrapeConfig.MetricConfig.Name

	desc := makeMetricDesc(config)

	assert(t, strings.Contains(desc.String(), "fqName: \""+expected), "expected metric name to be %s, got %s", expected, desc.String())
}

func TestMakeNewConstMetric_unsupportedMetricType(t *testing.T) {
	value := 123.0
	config := testExporterConfig
	config.ScrapeConfig.MetricConfig.Type = "summary"

	_, err := makeNewConstMetric(config, value)
	assert(t, err != nil, "makeNewConstMetric should return an error if the metric type is summary")
	errorContains(t, err, "not implemented yet")
}

func TestMakeNewConstMetric_errorCreatingMetric(t *testing.T) {
	value := 123.0
	config := testExporterConfig
	config.ScrapeConfig.MetricConfig.Name = "%invalid metric name!!%"

	_, err := makeNewConstMetric(config, value)
	assert(t, err != nil, "makeNewConstMetric should return an error for an invalid metric")
}

func TestGetLabelKeys(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
		"sun": "rain",
		"wet": "dry",
	}

	labelKeys := getLabelKeys(labels)
	expected := []string{"foo", "sun", "wet"}

	assert(t, compareStringSlices(labelKeys, expected), "expected to get the map keys %s. got: %s", expected, labelKeys)
}

func TestGetLabelValues(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
		"sun": "rain",
		"wet": "dry",
	}

	labelValues := getLabelValues(labels)
	expected := []string{"bar", "rain", "dry"}

	assert(t, compareStringSlices(labelValues, expected), "expected to get the map values %s. got: %s", expected, labelValues)
}

func TestGetPrometheusValueType(t *testing.T) {
	gaugeMetricType := getPrometheusValueType("gauge")

	assert(t, gaugeMetricType == prometheus.GaugeValue, "expected to return a GaugeValue for string 'gauge'. got enum with value: %d", gaugeMetricType)

	counterMetricType := getPrometheusValueType("counter")
	assert(t, counterMetricType == prometheus.CounterValue, "expected to return a CounterValue for string 'counter'. got enum with value: %d", counterMetricType)

	defaultMetricType := getPrometheusValueType("anything else")
	assert(t, defaultMetricType == prometheus.UntypedValue, "expected to return an UntypedValue for any other string. got enum with value: %d", defaultMetricType)
}
