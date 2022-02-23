package main

import (
	"fmt"
	"os"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

func getDefaultConfig() types.ExporterConfig {
	return types.ExporterConfig{
		GlobalConfig: types.GlobalConfig{
			MetricNamePrefix: "htmlexporter_",
			Port:             9883,
		},
	}
}

func getConfig(configFileArg *os.File) types.ExporterConfig {
	fileBytes, err := readConfigFile(configFileArg)
	if err != nil {
		fmt.Printf("error reading provided config file %s: %s", configFileArg.Name(), err)
		os.Exit(1)
	}

	config, err := parseConfig(fileBytes)
	if err != nil {
		fmt.Printf("error parsing provided config file %s: %s", configFileArg.Name(), err)
		os.Exit(1)
	}

	return config
}

func readConfigFile(file *os.File) ([]byte, error) {
	fileStat, err := file.Stat()
	if err != nil {
		// this error is very unlikely to occur, since `file` is already a descriptor to an open file
		return nil, fmt.Errorf("unable to stat file %s, invalid permissions or file does not exist. error: %s", file.Name(), err)
	}

	fileBytes := make([]byte, fileStat.Size())
	file.Read(fileBytes)

	return fileBytes, nil
}

func parseConfig(config []byte) (types.ExporterConfig, error) {
	exporterConfig := getDefaultConfig()

	err := yaml.UnmarshalStrict(config, &exporterConfig)
	if err != nil {
		return types.ExporterConfig{}, fmt.Errorf("error parsing supplied YAML configuration file: %s", err.Error())
	}

	addTargetDefaults(&exporterConfig)

	return exporterConfig, nil
}

func addTargetDefaults(config *types.ExporterConfig) {
	for i, target := range config.Targets {
		if target.DecimalPointSeparator == "" {
			config.Targets[i].DecimalPointSeparator = "."
		}

		if target.ThousandsSeparator == "" {
			config.Targets[i].ThousandsSeparator = ","
		}

		for j, metric := range target.Metrics {
			if metric.Type == "" {
				config.Targets[i].Metrics[j].Type = "untyped"
			}
		}
	}
}
