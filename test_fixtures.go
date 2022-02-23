package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var testExporterConfig = types.ExporterConfig{
	Targets: []types.TargetConfig{
		{
			DecimalPointSeparator: ".",
			ThousandsSeparator:    ",",
			Metrics: []types.MetricConfig{
				{
					Selector: "//div[@id='foobar']/text()",
					Name:     "wikipedia_articles_total",
					Type:     "gauge",
					Help:     "Total of articles available at Wikipedia",
					Labels: map[string]string{
						"language": "english",
					},
				},
			},
		},
		{
			DecimalPointSeparator: ",",
			ThousandsSeparator:    " ",
			Metrics: []types.MetricConfig{
				{
					Selector: "//div[@id='foobar']/text()",
					Name:     "wikipedia_articles_total",
					Type:     "gauge",
					Help:     "Total of articles available at Wikipedia",
					Labels: map[string]string{
						"language": "french",
					},
				},
			},
		},
	},
	GlobalConfig: types.GlobalConfig{
		MetricNamePrefix: "htmlexporter_",
	},
}

func getTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
}

type metricResult struct {
	name       string
	labels     map[string]string
	value      float64
	metricType dto.MetricType
}

// inspired from mysqld exporter test
// see: https://github.com/prometheus/mysqld_exporter/blob/e2ff660f50422245cdae9516dbf167e8c889c8bf/collector/collector_test.go#L31
func readMetric(m prometheus.Metric) metricResult {
	promMetric := &dto.Metric{}
	m.Write(promMetric)

	labels := make(map[string]string, len(promMetric.Label))

	for _, label := range promMetric.Label {
		labels[label.GetName()] = label.GetValue()
	}

	if promMetric.Gauge != nil {
		gauge := promMetric.GetGauge()
		return metricResult{
			name:       gauge.String(),
			value:      gauge.GetValue(),
			labels:     labels,
			metricType: dto.MetricType_GAUGE,
		}
	}

	if promMetric.Counter != nil {
		counter := promMetric.GetCounter()
		return metricResult{
			name:       counter.String(),
			value:      counter.GetValue(),
			labels:     labels,
			metricType: dto.MetricType_COUNTER,
		}
	}

	if promMetric.Untyped != nil {
		untyped := promMetric.GetUntyped()
		return metricResult{
			name:       untyped.String(),
			value:      untyped.GetValue(),
			labels:     labels,
			metricType: dto.MetricType_UNTYPED,
		}
	}

	panic(fmt.Sprintf("unsupported metric type %s", promMetric.String()))
}
