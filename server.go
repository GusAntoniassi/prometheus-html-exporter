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

type instrumentedResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	requestStart time.Time
}

func newInstrumentedResponseWriter(w http.ResponseWriter) *instrumentedResponseWriter {
	return &instrumentedResponseWriter{w, http.StatusOK, time.Now()}
}

func (irw *instrumentedResponseWriter) WriteHeader(code int) {
	irw.statusCode = code
	irw.ResponseWriter.WriteHeader(code)
}

func (irw *instrumentedResponseWriter) Header() http.Header {
	return irw.ResponseWriter.Header()
}

func (irw *instrumentedResponseWriter) Write(body []byte) (int, error) {
	return irw.ResponseWriter.Write(body)
}

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
	instrumentedResponseWriter := newInstrumentedResponseWriter(w)

	defer func() {
		duration := time.Since(instrumentedResponseWriter.requestStart).Seconds()

		// @TODO: expose metrics about duration and status

		log.WithFields(log.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"start":    instrumentedResponseWriter.requestStart.Local().UTC(),
			"duration": fmt.Sprintf("%0.3fs", duration),
			"code":     instrumentedResponseWriter.statusCode,
		}).Info()
	}()

	if !config.GlobalConfig.HasConfigFile {
		queryParams := r.URL.Query()
		if queryParams.Encode() == "" {
			instrumentedResponseWriter.WriteHeader(http.StatusNotFound)
			instrumentedResponseWriter.Write([]byte("no configuration supplied"))
			return
		}

	}

	collector := collector{config: config}
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(instrumentedResponseWriter, r)
}
