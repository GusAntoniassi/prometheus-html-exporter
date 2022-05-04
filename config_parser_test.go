package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
)

func openTestFile(t *testing.T, filename string) *os.File {
	testDir := getTestDir(t)

	sampleFile, err := os.Open(path.Join(testDir, filename))
	assert(t, err == nil, fmt.Sprintf("error opening sample file: %s. this is likely a problem in the test itself", err))

	return sampleFile
}

func TestReadConfigFile(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config.yaml")
	config, err := readConfigFile(sampleFile)

	ok(t, err)
	assert(t, len(config) > 0, "sample config file should not be empty")
}

func TestParseConfig(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config.yaml")
	configFile, err := readConfigFile(sampleFile)
	ok(t, err)

	config, err := parseConfig(configFile)
	ok(t, err)

	expected := "wikipedia_articles_total"
	actual := config.Targets[0].Metrics[0].Name

	assert(t, expected == actual, "metric name should be %s, got: '%s'", expected, actual)
}

func TestParseConfig_withDefaultValues(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config-minimal.yaml")
	configFile, err := readConfigFile(sampleFile)
	ok(t, err)

	config, err := parseConfig(configFile)
	ok(t, err)

	expectedActual := []map[string]string{
		{"desc": "metric_name_prefix", "expected": "htmlexporter_", "actual": config.GlobalConfig.MetricNamePrefix},
		{"desc": "port", "expected": "9883", "actual": fmt.Sprintf("%d", config.GlobalConfig.Port)},
	}

	for _, value := range expectedActual {
		assert(t, value["expected"] == value["actual"], "expected value of %s to be %s, got %s", value["desc"], value["expected"], value["actual"])
	}

	for _, target := range config.Targets {
		assert(t, target.DecimalPointSeparator == ".", "expected decimal point separator to be '.', got: '%s'", target.DecimalPointSeparator)
		assert(t, target.ThousandsSeparator == ",", "expected thousands separator to be ',', got: '%s'", target.ThousandsSeparator)
	}
}

func TestParseConfig_invalidParameters(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config-invalid-parameters.yaml")
	configFile, err := readConfigFile(sampleFile)
	ok(t, err)

	_, err = parseConfig(configFile)

	errorContains(t, err, "unmarshal errors")
}

func TestParseConfig_invalidYaml(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config-invalid-syntax.notyaml")
	configFile, err := readConfigFile(sampleFile)
	ok(t, err)

	_, err = parseConfig(configFile)

	errorContains(t, err, "error parsing supplied YAML")
}

func TestGetConfig(t *testing.T) {
	sampleFile := openTestFile(t, "sample-config.yaml")
	config := getConfig(sampleFile)

	assert(t, config.GlobalConfig.Port == 9883, "getConfig should contain the sample-config.yaml file configurations")
}

func TestGetConfig_invalidConfig(t *testing.T) {
	// getConfig may call os.Exit(1), so we have to get around that
	// see: https://talks.golang.org/2014/testing.slide#23
	// unfortunately since this is a subprocess we won't be able to get the full test coverage metric
	if os.Getenv("BE_CRASHER") == "1" {
		sampleFile := openTestFile(t, "sample-config-invalid-syntax.notyaml")
		getConfig(sampleFile)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestGetConfig_invalidConfig$")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok {
		assert(t, e.ExitCode() == 1, "process should exit with code 1, was: %d", e.ExitCode())
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestAddTargetDefaults(t *testing.T) {
	config := types.ExporterConfig{
		Targets: []types.TargetConfig{
			{Metrics: []types.MetricConfig{
				{},
			}},
		},
	}

	for i := range config.Targets {
		addTargetDefaults(&config.Targets[i])
	}

	assert(t, config.Targets[0].DecimalPointSeparator == ".", "expected default decimal point separator to be '.', got: '%s'", config.Targets[0].DecimalPointSeparator)
	assert(t, config.Targets[0].ThousandsSeparator == ",", "expected default thousands separator to be ',', got: '%s'", config.Targets[0].ThousandsSeparator)
	assert(t, config.Targets[0].Metrics[0].Type == "untyped", "expected default metric type to be 'untyped', got: '%s'", config.Targets[0].Metrics[0].Type)
}

func TestGetTargetConfigFromURLQuery(t *testing.T) {
	// @TODO
}
