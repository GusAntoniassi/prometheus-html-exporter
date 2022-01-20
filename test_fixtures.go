package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
)

var testExporterConfig = types.ExporterConfig{
	ScrapeConfig: types.ScrapeConfig{
		Selector:              "//div[@id='foobar']/text()",
		DecimalPointSeparator: ".",
		ThousandsSeparator:    ",",
		MetricConfig: types.MetricConfig{
			Name: "wikipedia_articles_total",
			Type: "gauge",
			Help: "Total of articles available at Wikipedia",
			Labels: map[string]string{
				"language": "english",
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
