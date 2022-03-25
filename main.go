package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	parser := argparse.NewParser("html-exporter", "Parses exported command-line configuration flags")

	configFile := parser.File("c", "config", os.O_RDONLY, 0600, &argparse.Options{
		Help:     "Path to the YAML configuration file",
		Required: false,
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	config := getDefaultConfig()
	if !argparse.IsNilFile(configFile) {
		config = getConfig(configFile)
	} else {
		log.Info("initializing exporter without a config file")
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
