package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

func openTestFile(t *testing.T, filename string) *os.File {
	testDir := getTestDir(t)

	sampleFile, err := os.Open(path.Join(testDir, filename))
	assert(t, err == nil, fmt.Sprintf("error opening sample file: %s. this is likely a problem in the test itself", err))

	return sampleFile
}

func TestGetDefaultConfig(t *testing.T) {
	config := getDefaultConfig()

	assert(t, config.GlobalConfig.Port == 9883, "default port should be 9883, got %d", config.GlobalConfig.Port)
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

	assert(t, config.ScrapeConfig.MetricConfig.Name == "wikipedia_articles_total", "metric name should be 'wikipedia_articles_total', got: %s", config.ScrapeConfig.MetricConfig.Name)
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
