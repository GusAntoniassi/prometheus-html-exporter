package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {

	config := types.ExporterConfig{
		ScrapeConfig: types.ScrapeConfig{
			Address:               "https://en.wikipedia.org/wiki/Special:Statistics",
			Selector:              "//div[@id='mw-content-text']//tr[@class='mw-statistics-articles']/td[@class='mw-statistics-numbers']/text()",
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
			Port:             9883,
		},
	}

	metricRegistry, err := getExporterMetricsRegistry()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, config)
	})

	http.Handle("/metrics", promhttp.HandlerFor(metricRegistry, promhttp.HandlerOpts{}))

	server := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", config.GlobalConfig.Port),
	}

	log.Infof("Server starting and listening on port %d", config.GlobalConfig.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %s", err.Error())
	}
}
