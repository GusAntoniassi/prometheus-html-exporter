package main

import (
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
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

func TestNormalizeNumericValue_nbspSpacedValue(t *testing.T) {
	normalizedValue, err := normalizeNumericValue("1\u00a0234\u00a0567.08", " ", ".")
	ok(t, err)
	assert(t, normalizedValue == expectedNormalizedValue, "expected normalized value with nbsp and space separators to equal to 1234567.08. actual value was %f", normalizedValue)
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

func TestParseDomContent(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("<html></html>"))
	document, err := parseDomContent(reader)

	ok(t, err)
	assert(t, document.FirstChild.Data == "html", "first element of DOM document should be 'html'. got: %s", document.FirstChild.Data)
}

func TestParseSelector(t *testing.T) {
	expected := "Hello world"
	reader := io.NopCloser(strings.NewReader("<html><body><div id=\"foobar\">Hello world</div></body></html>"))
	document, err := parseDomContent(reader)
	ok(t, err)

	output, err := parseSelector(document, "//div[@id='foobar']/text()")
	ok(t, err)

	assert(t, output == expected, "expected \"Hello world\" text selected by the XPath expression, got: %s", output)
}

func TestParseSelector_invalidXPath(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("<html></html>"))
	document, err := parseDomContent(reader)
	ok(t, err)

	_, err = parseSelector(document, "/`$/")
	assert(t, err != nil, "expected error for an invalid XPath expression")

	errorContains(t, err, "querying the XPath")
}

func TestParseSelector_emptyElements(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("<html></html>"))
	document, err := parseDomContent(reader)
	ok(t, err)

	_, err = parseSelector(document, "//div")
	assert(t, err != nil, "expected error when no elements were returned by the XPath query")
}

func TestParseSelector_warnsMoreThanOneElement(t *testing.T) {
	buf := captureLogOutput()

	reader := io.NopCloser(strings.NewReader("<html><div></div><div></div></html>"))
	document, err := parseDomContent(reader)
	ok(t, err)

	parseSelector(document, "//div")

	logOutput := buf.String()
	assertWarningLog(t, logOutput)
}

func TestScrape(t *testing.T) {
	testConfig := []struct {
		html     string
		expected float64
		server   *httptest.Server
	}{
		{
			html:     "<html><div id=\"foo\">1,234,567.08</div><div id=\"bar\">1,234,567.08</div></html>",
			expected: 1234567.08,
		},
		{
			html:     "<div id=\"bar\">987.654.321,00</div>",
			expected: 987654321.00,
		},
	}

	testConfig[0].server = getTestServer(testConfig[0].html)
	testConfig[1].server = getTestServer(testConfig[1].html)

	config := []types.TargetConfig{
		{
			Address:               testConfig[0].server.URL,
			DecimalPointSeparator: ".",
			ThousandsSeparator:    ",",
			Metrics: []types.MetricConfig{
				{
					Name:     "foo",
					Selector: "//div[@id='foo']/text()",
				},
				{
					Name:     "bar",
					Selector: "//div[@id='foo']/text()",
				},
			},
		},
		{
			Address:               testConfig[1].server.URL,
			DecimalPointSeparator: ",",
			ThousandsSeparator:    ".",
			Metrics: []types.MetricConfig{
				{
					Selector: "//div[@id='bar']/text()",
				},
			},
		},
	}

	output := scrape(config)

	assert(t, output[0][0] == testConfig[0].expected, "expected first scrape / first value to be equal to %0.2f, got %0.2f", testConfig[0].expected, output[0][0])
	assert(t, output[0][1] == testConfig[0].expected, "expected first scrape / second value to be equal to %0.2f, got %0.2f", testConfig[0].expected, output[0][1])
	assert(t, output[1][0] == testConfig[1].expected, "expected second scrape / first value to be equal to %0.2f, got %0.2f", testConfig[1].expected, output[1][0])
}

var testTargetConfig = []types.TargetConfig{
	{
		Address:               "https://en.wikipedia.org/wiki/Special:Statistics",
		DecimalPointSeparator: ".",
		ThousandsSeparator:    ",",
		Metrics: []types.MetricConfig{
			{
				Selector: "//div[@id='foobar']/text()",
			},
		},
	},
}

func TestScrape_requestError(t *testing.T) {
	buf := captureLogOutput()

	config := testTargetConfig
	config[0].Address = "http://go%Qdev"

	scrape(config)

	logOutput := buf.String()
	logContainsScrapeError := strings.Contains(logOutput, "error scraping target")

	assertWarningLog(t, logOutput)
	assert(t, logContainsScrapeError, "expected warning log when there is an error scraping a target.")
}

func TestScrapeTarget(t *testing.T) {
	html := "<div id=\"foobar\">1,234,567.08</div>"
	expected := 1234567.08

	server := getTestServer(html)

	config := testTargetConfig[0]
	config.Address = server.URL

	output, err := scrapeTarget(config)

	ok(t, err)
	assert(t, output[0] == expected, "expected scrape value to be equal to %0.2f, got %0.2f", expected, output)
}

func TestScrapeTarget_requestError(t *testing.T) {
	config := testTargetConfig[0]
	config.Address = "http://go%Qdev"

	_, err := scrapeTarget(config)
	assert(t, err != nil, "expected scrape to return an error when the HTTP request fails")
}

func TestScrapeTarget_warnsInvalidXPath(t *testing.T) {
	buf := captureLogOutput()

	html := ""
	xpath := "/^$/"

	server := getTestServer(html)

	config := testTargetConfig[0]
	config.Metrics[0].Selector = xpath
	config.Address = server.URL

	values, err := scrapeTarget(config)
	ok(t, err)

	logOutput := buf.String()
	isXPathWarningLog := strings.Contains(logOutput, "with XPATH selector")

	assertWarningLog(t, logOutput)
	assert(t, isXPathWarningLog, "expected warning about XPATH expression error. log output was: %s", logOutput)

	assert(t, math.IsNaN(values[0]), "expected errored metric value to be NaN")
}

func TestScrape_warnOnNormalizeDomElement(t *testing.T) {
	buf := captureLogOutput()

	html := "<div id=\"foobar\">1,234,567.08</div>"
	// this xpath expression has no `/text()` function, so it will return element name ("div")
	xpath := "//div[@id='foobar']"

	server := getTestServer(html)

	config := testTargetConfig[0]
	config.Metrics[0].Selector = xpath
	config.Address = server.URL

	values, err := scrapeTarget(config)
	ok(t, err)

	logOutput := buf.String()
	isNormalizingWarningLog := strings.Contains(logOutput, "error normalizing value")

	assertWarningLog(t, logOutput)
	assert(t, isNormalizingWarningLog, "expected warning about normalizing non-numeric value. log output was: %s", logOutput)

	assert(t, math.IsNaN(values[0]), "expected errored metric value to be NaN")
}
