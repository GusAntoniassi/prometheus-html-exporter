package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func getExporterMetricsRegistry() (*prometheus.Registry, error) {
	metricRegistry := prometheus.NewRegistry()

	err := metricRegistry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	if err != nil {
		return nil, fmt.Errorf("error registering process collector: %s", err.Error())
	}

	err = metricRegistry.Register(collectors.NewGoCollector())
	if err != nil {
		return nil, fmt.Errorf("error registering Go collector: %s", err.Error())
	}

	return metricRegistry, nil
}

func probeHandler(w http.ResponseWriter, r *http.Request, config types.ExporterConfig) {
	// @TODO: gather some configs from query parameters, passed from Prometheus
	start := time.Now()

	collector := collector{config: config}
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

	duration := time.Since(start).Seconds()
	// @TODO: expose metrics about duration
	log.Debugf("scrape of all targets finished in %0.2f seconds", duration)
}
