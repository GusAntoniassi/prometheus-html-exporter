# Documentation

## Motivation
This exporter aims to facilitate the implementation of simple web scraping, without requiring you to rollout your own code, and also allowing you to use Prometheus' statistics queries and features on the scraped values.

For instance, let's say you needed a metric on how many Wikipedia articles there are. There is a [page for that](https://en.wikipedia.org/wiki/Special:Statistics), and the value you would need is inside the element `<td class="mw-statistics-numbers">6,440,382</td>`. You could write a bot that scrapes that page using the following XPath selector: `//td[@class='mw-statistics-numbers']/text()`.

Prometheus HTML exporter does it automatically for you, so you wouldn't have to write any code, just a simple YAML config:

```yaml
targets:
- address: https://en.wikipedia.org/wiki/Special:Statistics
  metrics:
  - name: wikipedia_articles_total
    type: gauge
    help: Total of articles available at Wikipedia
    selector: //div[@id='mw-content-text']//tr[@class='mw-statistics-articles']/td[@class='mw-statistics-numbers']/text()
```

## Configuring
Work in progress

### Necessary tools:

Most dependencies are managed by `asdf`, install them running:
```
asdf install
```

To template markdown files, use [`emd`](https://github.com/mh-cbon/emd).

## Developing
Work in progress

### Release process
Test everything before publishing:
```sh
make test
```

Update documentation (`README.md` and/or `docs/*.e.md`) if necessary:
```sh
make docs
```

Build binaries and Dockerfiles:
```sh
make release
```

Update CHANGELOG.md with current version info:
```sh
make changelog
git commit -m "chore: update CHANGELOG" ./CHANGELOG.md
```

Generate version tag:
```sh
git tag vX.X.X
git push --tags
```

[Create the release](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository) for the current tag, and upload files in `build/tarball/*` to GitHub

Push the image to Docker Hub:
```sh
docker push gusantoniassi/prometheus-html-exporter
```
