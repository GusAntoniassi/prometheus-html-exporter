# Prometheus HTML exporter

**Still under development**

This project aims to allow the scraping of metrics contained within HTML web pages, collecting them with an XPath selector and exporting them as metrics in a Prometheus format. See the [docs](docs/documentation.md) for more information.

## Running (under development)
For now, you must pass a config yaml file, so you could run the program as:
```sh
go run . -c examples/full-config.yaml
```
A binary release distribution and Docker image are planned for the near future.

### Testing

Run the test suite with the following command:
```sh
make test
```

## Features
### Implemented
- Scrape a web value using XPath
- YAML file configuration

### Under development:
- Binary and Docker image releases
- Query param configuration (allows native integration with Prometheus `scrape_configs`)
- Multiple endpoint configuration
- Exporter instrumentation (metrics about the scrape itself)
- Timeouts
- Basic auth scrape
- Basic arithmetic with scraped value
- Arithmetic "pipeline" for one or more scraped values (e.g. allowing you to divide two numbers)
