package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func getTestDir(tb testing.TB) string {
	cwd, err := os.Getwd()
	assert(tb, err == nil, fmt.Sprintf("error getting cwd: %s. this is likely a problem in the test itself", err))

	return path.Join(cwd, "testdata")
}

// Taken from https://github.com/benbjohnson/testing

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Fatalf(msg, v...)
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("unexpected error: %s", err.Error())
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Fatalf("exp: %#v\n\n\tgot: %#v", exp, act)
	}
}

func errorContains(tb testing.TB, err error, errorSubstr string) {
	assert(tb, err != nil, "err should not be nil")
	assert(tb, strings.Contains(err.Error(), errorSubstr), "expected error containing %q, got: %s", errorSubstr, err.Error())
}

func assertWarningLog(tb testing.TB, logOutput string) {
	assert(tb, strings.Contains(logOutput, "\"level\":\"warning\""), "expected a warning log message. log output was: %s", logOutput)
}

func compareStringSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// since go map ordering is random we can't rely on it to compare the values
	// https://nathanleclaire.com/blog/2014/04/27/a-surprising-feature-of-golang-that-colored-me-impressed
	for _, valueA := range a {
		match := false

		for _, valueB := range b {
			if valueA == valueB {
				match = true
			}
		}

		if !match {
			return false
		}
	}

	return true
}

func captureLogOutput() *bytes.Buffer {
	log.SetFormatter(&log.JSONFormatter{})

	var buf bytes.Buffer
	log.SetOutput(&buf)

	return &buf
}
