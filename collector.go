package main

import (
	"fmt"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	config types.ExporterConfig
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	for _, target := range c.config.Targets {
		for _, metric := range target.Metrics {
			ch <- makeMetricDesc(c.config, metric)
		}
	}
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	values := scrape(c.config.Targets)

	for i, target := range c.config.Targets {
		for j, metric := range target.Metrics {
			prometheusMetric, err := makeNewConstMetric(c.config, metric, values[i][j])

			if err != nil {
				panic(fmt.Sprintf("error making const metric for %s: %s", metric.Name, err.Error()))
			}

			ch <- prometheusMetric
		}
	}
}

func makeMetricDesc(config types.ExporterConfig, metric types.MetricConfig) *prometheus.Desc {
	return prometheus.NewDesc(
		config.GlobalConfig.MetricNamePrefix+metric.Name,
		metric.Help,
		getLabelKeys(metric.Labels),
		nil,
	)
}

func makeNewConstMetric(config types.ExporterConfig, metric types.MetricConfig, value float64) (prometheus.Metric, error) {
	var valueType prometheus.ValueType

	switch metric.Type {
	case "histogram":
	case "summary":
		return nil, fmt.Errorf("metric type \"%s\" is not supported", metric.Type)
	default:
		valueType = getPrometheusValueType(metric.Type)
	}

	desc := makeMetricDesc(config, metric)

	labelValues := getLabelValues(metric.Labels)
	prometheusMetric, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)

	if err != nil {
		return nil, err
	}

	return prometheusMetric, nil
}

func getLabelKeys(labels map[string]string) []string {
	labelKeys := make([]string, len(labels))

	i := 0
	for k := range labels {
		labelKeys[i] = k
		i++
	}

	return labelKeys
}

func getLabelValues(labels map[string]string) []string {
	labelValues := make([]string, len(labels))

	i := 0
	for _, v := range labels {
		labelValues[i] = v
		i++
	}

	return labelValues
}

func getPrometheusValueType(metricType string) prometheus.ValueType {
	var valueType prometheus.ValueType

	if metricType == "gauge" {
		valueType = prometheus.GaugeValue
	} else if metricType == "counter" {
		valueType = prometheus.CounterValue
	} else {
		valueType = prometheus.UntypedValue
	}

	return valueType
}
