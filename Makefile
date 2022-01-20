.PHONY: run
run:
	go run . -c examples/full-config.yaml

.PHONY: test
test:
	go test -v
