package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCollect(t *testing.T) {
	config := testExporterConfig
	expectedLabels := make([]map[string]string, 0, len(config.Targets))

	for i, target := range config.Targets {
		html := fmt.Sprintf("<div id=\"foobar\">1%s234%s567%s08</div>", target.ThousandsSeparator, target.ThousandsSeparator, target.DecimalPointSeparator)

		server := getTestServer(html)

		config.Targets[i].Address = server.URL

		for _, metric := range target.Metrics {
			expectedLabels = append(expectedLabels, metric.Labels)
		}
	}

	collector := collector{config: config}

	ch := make(chan prometheus.Metric)

	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	counter := 0
	for metric := range ch {
		assert(t, metric != nil, "collect should not return a nil metric")

		result := readMetric(metric)
		assert(t, result.value == 1234567.08, "expected scraped metric value to be 1234567.08, got %0.2f", result.value)
		assert(t, reflect.DeepEqual(result.labels, expectedLabels[counter]), "expected scraped metric label to be %#v, got %#v", expectedLabels[counter], result.labels)

		counter++
	}
}

func TestCollect_errorMakingConstMetric(t *testing.T) {
	server := getTestServer("<div id=\"foobar\">123</div>")
	config := testExporterConfig
	config.Targets = []types.TargetConfig{{
		Address: server.URL,
		Metrics: []types.MetricConfig{{
			Name: "foo",
			Type: "summary",
		}},
	}}

	collector := collector{config: config}

	// recover from panic
	defer func() { recover() }()

	ch := make(chan prometheus.Metric)

	collector.Collect(ch)
	t.Errorf("should panic when configured with a summary metric")
}

func TestDescribe(t *testing.T) {
	config := testExporterConfig
	config.Targets = []types.TargetConfig{
		{Metrics: []types.MetricConfig{
			{Name: "foo"},
			{Name: "bar"},
		}},
		{Metrics: []types.MetricConfig{
			{Name: "foobar"},
		}},
	}

	prefix := config.GlobalConfig.MetricNamePrefix
	expectedMetricNames := []string{
		prefix + "foo",
		prefix + "bar",
		prefix + "foobar",
	}

	collector := collector{config: config}

	ch := make(chan *prometheus.Desc)

	go func() {
		collector.Describe(ch)
		close(ch)
	}()

	counter := 0
	for desc := range ch {
		assert(
			t,
			strings.Contains(
				desc.String(),
				"fqName: \""+expectedMetricNames[counter],
			),
			"expected metric name to be %s, got %s",
			expectedMetricNames[counter],
			desc.String(),
		)
		counter++
	}
}

func TestCollect_scrapeError(t *testing.T) {
	config := testExporterConfig
	config.Targets[0].Address = "foo://bar.dev"

	collector := collector{config: config}

	// recover from panic
	defer func() { recover() }()

	ch := make(chan prometheus.Metric, 1)

	collector.Collect(ch)
	t.Errorf("should panic when configured with a summary metric")
}

func TestMakeMetricDesc(t *testing.T) {
	config := testExporterConfig

	for _, target := range config.Targets {
		for _, metric := range target.Metrics {
			expected := config.GlobalConfig.MetricNamePrefix + metric.Name

			desc := makeMetricDesc(config, metric)

			assert(t, strings.Contains(desc.String(), "fqName: \""+expected), "expected metric name to be %s, got %s", expected, desc.String())
		}
	}
}

func TestMakeNewConstMetric(t *testing.T) {
	value := 123.0
	metric := types.MetricConfig{
		Name: "foobar",
		Type: "gauge",
	}
	_, err := makeNewConstMetric(testExporterConfig, metric, value)

	ok(t, err)
}

func TestMakeNewConstMetric_unsupportedMetricType(t *testing.T) {
	value := 123.0
	config := testExporterConfig
	metric := types.MetricConfig{
		Name: "foobar",
		Type: "summary",
	}

	_, err := makeNewConstMetric(config, metric, value)
	assert(t, err != nil, "makeNewConstMetric should return an error if the metric type is summary")
	errorContains(t, err, "not supported")
}

func TestMakeNewConstMetric_errorCreatingMetric(t *testing.T) {
	value := 123.0
	config := testExporterConfig
	metric := types.MetricConfig{
		Name: "%invalid metric name!!%",
	}

	_, err := makeNewConstMetric(config, metric, value)
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
