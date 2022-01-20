package main

import (
	"fmt"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type collector struct {
	config types.ExporterConfig
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	log.Info("begin collect")

	value, err := scrape(c.config.ScrapeConfig)

	if err != nil {
		// @TODO: better handling
		panic(fmt.Sprintf("error scraping: %s", err.Error()))
	}

	metric, err := makeNewConstMetric(c.config, value)

	if err != nil {
		panic(err.Error())
	}

	ch <- metric
}

func makeNewConstMetric(config types.ExporterConfig, value float64) (prometheus.Metric, error) {
	metricConfig := config.ScrapeConfig.MetricConfig
	var valueType prometheus.ValueType

	switch metricConfig.Type {
	case "summary":
		return nil, fmt.Errorf("not implemented yet")
	case "histogram":
		return nil, fmt.Errorf("not implemented yet")
	default:
		valueType = getPrometheusValueType(metricConfig.Type)
	}

	desc := prometheus.NewDesc(
		config.GlobalConfig.MetricNamePrefix+metricConfig.Name,
		metricConfig.Help,
		getLabelKeys(metricConfig.Labels),
		nil,
	)

	labelValues := getLabelValues(metricConfig.Labels)
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)

	if err != nil {
		return nil, err
	}

	return metric, nil
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
