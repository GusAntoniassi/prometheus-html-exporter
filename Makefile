.PHONY: run
run:
	go run . -c examples/full-config.yaml

.PHONY: test
test:
	go test -v

.PHONY: README.md
README.md:
	emd gen -in README.e.md -out README.md

.PHONY: docs/documentation.md
docs/documentation.md:
	emd gen -in docs/documentation.e.md -out docs/documentation.md

.PHONY: docs
docs: README.md docs/documentation.md

.PHONY: changelog
changelog:
	changelog prepare
	vim change.log
