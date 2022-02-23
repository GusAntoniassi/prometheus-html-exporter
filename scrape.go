package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/GusAntoniassi/prometheus-html-exporter/internal/pkg/types"
	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func scrape(targets []types.TargetConfig) [][]float64 {
	targetCount := len(targets)
	scrapedValues := make([][]float64, targetCount)

	for i, target := range targets {
		targetValues, err := scrapeTarget(target)

		// An error here will only mean there is an issue with the request itself
		// Failed individual metrics will return NaN if there is an error
		if err != nil {
			log.Warnf("error scraping target %s: %s", target.Address, err)
			// @TODO: increment scrape error metric
			continue
		}

		scrapedValues[i] = targetValues
	}

	return scrapedValues
}

func scrapeTarget(target types.TargetConfig) ([]float64, error) {
	metricCount := len(target.Metrics)
	scrapedValues := make([]float64, metricCount)

	// initialize the slice with NaN to allow early returns and identify errors later on
	for i := range scrapedValues {
		scrapedValues[i] = math.NaN()
	}

	log.Debugf("requesting URL '%s'", target.Address)
	body, err := doRequest(target.Address)
	if err != nil {
		return scrapedValues, err
	}

	document, err := parseDomContent(body)
	if err != nil {
		return scrapedValues, err
	}

	for i, metric := range target.Metrics {
		log.Debugf("scraping value from requested URL with XPath selector '%s'", metric.Selector)
		scrapedValue, err := parseSelector(document, metric.Selector)

		if err != nil {
			log.Warnf("error scraping value for metric %s with XPATH selector '%s'. error: %s", metric.Name, metric.Selector, err.Error())
			continue
		}

		numberValue, err := normalizeNumericValue(scrapedValue, target.ThousandsSeparator, target.DecimalPointSeparator)
		if err != nil {
			log.Warnf("error normalizing value %s for metric %s with XPATH selector '%s'. error: %s", scrapedValue, metric.Name, metric.Selector, err.Error())
			continue
		}

		log.Debugf("scraped value '%0.2f' from URL '%s'", numberValue, target.Address)
		scrapedValues[i] = numberValue
	}

	return scrapedValues, nil
}

func doRequest(url string) (io.ReadCloser, error) {
	// @TODO: Allow passing headers, timeout and other request args
	client := &http.Client{
		// Timeout: 10000,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request. error: %s", err)
	}

	req.Header.Add("User-Agent", fmt.Sprintf("prometheus-html-exporter/%s", BuildVersion))

	log.Infof("scraping page %s", url)

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("unable to request URL %s. error: %s", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 400 {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	return resp.Body, nil
}

func parseDomContent(body io.ReadCloser) (*html.Node, error) {
	doc, err := htmlquery.Parse(body)

	if err != nil {
		// this error is mostly relating to body reading errors so it won't be covered in tests
		return nil, fmt.Errorf("error loading the response body into XPath nodes. error: %s", err)
	}

	return doc, nil
}

func parseSelector(document *html.Node, selector string) (string, error) {
	nodes, err := htmlquery.QueryAll(document, selector)

	if err != nil {
		return "", fmt.Errorf("error querying the XPath expression `%s`. error: %s", selector, err)
	}

	if len(nodes) < 1 {
		return "", fmt.Errorf("no elements returned by the XPath expression `%s`", selector)
	}

	// currently supporting only one attribute. this could change in the future if necessary
	if len(nodes) > 1 {
		log.Warn("more than one element was returned by the XPath expression. only the value of the first element will be exported")
	}

	value := nodes[0].Data

	return value, nil
}

func normalizeNumericValue(value string, thousandsSeparator string, decimalSeparator string) (float64, error) {
	// Replace non-breaking space characters with regular space before parsing
	value = strings.ReplaceAll(value, "\u00a0", " ")

	// Replace separators to convert the string into a format accepted by strconv
	value = strings.ReplaceAll(strings.ReplaceAll(value, thousandsSeparator, ""), decimalSeparator, ".")

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing value %s to a float. error: %s", value, err)
	}

	return floatValue, nil
}
