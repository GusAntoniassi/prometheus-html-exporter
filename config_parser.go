package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

func getDefaultConfig() types.ExporterConfig {
	return types.ExporterConfig{
		GlobalConfig: types.GlobalConfig{
			MetricNamePrefix: "htmlexporter_",
			Port:             9883,
			HasConfigFile:    false, // will be overriden
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

	config.GlobalConfig.HasConfigFile = true

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

	for i := range exporterConfig.Targets {
		addTargetDefaults(&exporterConfig.Targets[i])
	}

	return exporterConfig, nil
}

func getTargetConfigFromURLQuery(query url.Values) ([]types.TargetConfig, error) {
	target := types.TargetConfig{
		Address:               query.Get("target"),
		DecimalPointSeparator: query.Get("decimal_point_separator"),
		ThousandsSeparator:    query.Get("thousands_separator"),
		Metrics: []types.MetricConfig{
			{
				Name:     query.Get("metric_name"),
				Help:     query.Get("metric_help"),
				Type:     query.Get("metric_type"),
				Selector: query.Get("selector"),
			},
		},
	}

	addTargetDefaults(&target)

	err := types.Validate(target)

	if err != nil {
		return []types.TargetConfig{}, err
	}

	return []types.TargetConfig{target}, nil
}

func addTargetDefaults(target *types.TargetConfig) {
	if target.DecimalPointSeparator == "" {
		target.DecimalPointSeparator = "."
	}

	if target.ThousandsSeparator == "" {
		target.ThousandsSeparator = ","
	}

	for j, metric := range target.Metrics {
		if metric.Type == "" {
			target.Metrics[j].Type = "untyped"
		}
	}
}
