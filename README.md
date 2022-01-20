# Prometheus HTML exporter

**Still under development**

This project aims to allow the scraping of metrics contained within HTML web pages, collecting them with an XPath selector and exporting them as metrics in a Prometheus format. See the [docs](docs/documentation.md) for more information.

## Running (under development)
For now, the configuration stays in the `main.go` file and starting the server should be done with the command:
```sh
make run
```

## Development
### Testing

Run the test suite with the following command:
```sh
make test
```

## Features
### Implemented
- Scrape a web value using XPath

### Under development:
- YAML file configuration
- Query param configuration (allows native integration with Prometheus `scrape_configs`)
- Multiple endpoint configuration
- Exporter instrumentation (metrics about the scrape itself)
- Timeouts
- Basic auth scrape
- Basic arithmetic with scraped value
