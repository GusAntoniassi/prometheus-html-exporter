# Documentation

## Motivation
This exporter aims to facilitate the implementation of simple web scraping, without requiring you to rollout your own code, and also allowing you to use Prometheus' statistics queries and features on the scraped values.

For instance, let's say you needed a metric on how many Wikipedia articles there are. There is a [page for that](https://en.wikipedia.org/wiki/Special:Statistics), and the value you would need is inside the element `<td class="mw-statistics-numbers">6,440,382</td>`. You could write a bot that scrapes that page using the following XPath selector: `//td[@class='mw-statistics-numbers']/text()`.

Prometheus HTML exporter does it automatically for you, so you wouldn't have to write any code, just a simple YAML config:

```yaml
targets:
  - address: "https://en.wikipedia.org/wiki/Special:Statistics"
    metrics:
      - name: wikipedia_articles_total
        type: gauge
        help: "Total of articles available at Wikipedia"
        selector: "//div[@id='mw-content-text']//tr[@class='mw-statistics-articles']/td[@class='mw-statistics-numbers']/text()"
        labels:
          language: english
```

## Configuring
Work in progress

## Developing
Work in progress
