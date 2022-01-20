package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	log "github.com/sirupsen/logrus"
)

var expectedNormalizedValue = 1234567.08

func TestNormalizeNumericValue(t *testing.T) {
	enValue, err := normalizeNumericValue("1,234,567.08", ",", ".")
	ok(t, err)
	assert(t, enValue == expectedNormalizedValue, "expected value to equal to 1234567.08. actual value was %f", enValue)
}

func TestNormalizeNumericValue_noCommaValue(t *testing.T) {
	enValueNoComma, err := normalizeNumericValue("1234567.08", ",", ".")
	ok(t, err)
	assert(t, enValueNoComma == expectedNormalizedValue, "expected value to equal to 1234567.08. actual value was %f", enValueNoComma)
}

func TestNormalizeNumericValue_brValue(t *testing.T) {
	brValue, err := normalizeNumericValue("1.234.567,08", ".", ",")
	ok(t, err)
	assert(t, brValue == expectedNormalizedValue, "expected value to equal to 1234567.08. actual value was %f", brValue)
}

func TestNormalizeNumericValue_frValue(t *testing.T) {
	frValue, err := normalizeNumericValue("1 234 567,08", " ", ",")
	ok(t, err)
	assert(t, frValue == expectedNormalizedValue, "expected value to equal to 1234567.08. actual value was %f", frValue)
}

func TestNormalizeNumericValue_invalidValue(t *testing.T) {
	number, err := normalizeNumericValue("1234@@567!08", ",", ".")

	assert(t, err != nil, "expected normalizing an invalid numeric value to return an error. number was %0.2f", number)
}

func TestDoRequest(t *testing.T) {
	response := "<html>foobar</html>"

	server := getTestServer(response)

	output, err := doRequest(server.URL)
	ok(t, err)

	buffer, err := io.ReadAll(output)
	ok(t, err)

	outputString := strings.TrimSpace(string(buffer))

	assert(t, response == outputString, "expected server response body output %s, got: %s", response, outputString)
}

func TestDoRequest_invalidRequest(t *testing.T) {
	// invalid URL escaping makes http.NewRequest's validation to fail
	_, err := doRequest("http://go%Qdev")
	assert(t, err != nil, "expected doRequest to return an error on an invalid URL")
}

func TestDoRequest_erroredRequest(t *testing.T) {
	// Since the only way for client.Do to return an error is with a client policy
	// we return a response with a redirect to an invalid URL,
	// causing the `client.Do` call to fail
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "foobar://go.dev", http.StatusTemporaryRedirect)
	}))

	_, err := doRequest(server.URL)
	assert(t, err != nil, "expected doRequest to return an error when the request fails")
}

func TestDoRequest_erroredResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server error :(", 500)
	}))

	_, err := doRequest(server.URL)
	assert(t, err != nil, "expected doRequest to return an error when the server responds with an error")
}

func TestParseSelector(t *testing.T) {
	expected := "Hello world"
	reader := io.NopCloser(strings.NewReader("<html><body><div id=\"foobar\">Hello world</div></body></html>"))

	output, err := parseSelector(reader, "//div[@id='foobar']/text()")
	ok(t, err)

	assert(t, output == expected, "expected \"Hello world\" text selected by the XPath expression, got: %s", output)
}

func TestParseSelector_invalidXPath(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("<html></html>"))

	_, err := parseSelector(reader, "/`$/")
	assert(t, err != nil, "expected error for an invalid XPath expression")

	errorContains(t, err, "querying the XPath")
}

func TestParseSelector_emptyElements(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("<html></html>"))

	_, err := parseSelector(reader, "//div")
	assert(t, err != nil, "expected error when no elements were returned by the XPath query")
}

func TestParseSelector_warnsMoreThanOneElement(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})

	var buf bytes.Buffer
	log.SetOutput(&buf)

	reader := io.NopCloser(strings.NewReader("<html><div></div><div></div></html>"))
	_, _ = parseSelector(reader, "//div")

	logOutput := buf.String()
	isWarningLog := strings.Contains(logOutput, "\"level\":\"warning\"")
	assert(t, isWarningLog, "expected warning log when more than one element is returned by the XPath expression. log output was: %s", logOutput)
}

var testScrapeConfig = types.ScrapeConfig{
	Address:               "https://en.wikipedia.org/wiki/Special:Statistics",
	Selector:              "//div[@id='foobar']/text()",
	DecimalPointSeparator: ".",
	ThousandsSeparator:    ",",
}

func TestScrape(t *testing.T) {
	html := "<div id=\"foobar\">1,234,567.08</div>"
	expected := 1234567.08

	server := getTestServer(html)

	config := testScrapeConfig
	config.Address = server.URL

	output, err := scrape(config)
	ok(t, err)

	assert(t, output == expected, "expected scrape value to be equal to %0.2f, got %0.2f", expected, output)
}

func TestScrape_errorOnDivToNumber(t *testing.T) {
	html := "<div id=\"foobar\">1,234,567.08</div>"
	// this xpath expression has no `/text()` function, so it will return element name ("div")
	xpath := "//div[@id='foobar']"

	server := getTestServer(html)

	config := testScrapeConfig
	config.Selector = xpath
	config.Address = server.URL

	_, err := scrape(config)
	errorContains(t, err, "parsing value")
}

func TestScrape_invalidRequest(t *testing.T) {
	// invalid URL escaping makes http.NewRequest's validation to fail
	config := testScrapeConfig
	config.Address = "http://go%Qdev"

	_, err := scrape(config)
	assert(t, err != nil, "expected scrape to return an error when the HTTP request fails")
}

func TestScrape_invalidXPath(t *testing.T) {
	html := ""
	xpath := "/^$/"

	server := getTestServer(html)

	config := testScrapeConfig
	config.Selector = xpath
	config.Address = server.URL

	_, err := scrape(config)
	errorContains(t, err, "querying the XPath")
}
