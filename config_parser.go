package main

import (
	"fmt"
	"os"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

func getDefaultConfig() types.ExporterConfig {
	return types.ExporterConfig{
		ScrapeConfig: types.ScrapeConfig{
			DecimalPointSeparator: ".",
			ThousandsSeparator:    ",",
		},
		GlobalConfig: types.GlobalConfig{
			MetricNamePrefix: "htmlexporter_",
			Port:             9883,
		},
	}
}

func getConfig(configFileArg *os.File) types.ExporterConfig {
	fileBytes, err := readConfigFile(configFileArg)
	if err != nil {
		fmt.Printf("error reading config file: %s", err)
		os.Exit(1)
	}

	config, err := parseConfig(fileBytes)
	if err != nil {
		fmt.Printf("error parsing config file: %s", err)
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

	return exporterConfig, nil
}
