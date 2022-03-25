ROOT_DIR = $(shell pwd)
EXAMPLES_DIR = $(ROOT_DIR)/examples
BUILD_DIR = $(ROOT_DIR)/build
BIN_DIR = $(BUILD_DIR)/bin
TARBALL_DIR = $(BUILD_DIR)/tarball

# If GOOS is not specified, compile for all available OSs
ifeq ($(GOOS),)
GOOS := darwin dragonfly freebsd linux netbsd openbsd solaris windows
endif

DOCKER_IMAGE="gusantoniassi/prometheus-html-exporter"

.PHONY: dependency
dependency:
	go mod download
	asdf install
	pre-commit install

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
	git-chglog --output CHANGELOG.md
	vim CHANGELOG.md

.PHONY: build
build:
	for os in $(GOOS); do \
		os_arch="$(GOARCH)"; \
		if [ -z "$$os_arch" ]; then \
			case $$os in \
				darwin) \
					os_arch="amd64 arm64"; \
					;; \
				dragonfly) \
					os_arch="amd64"; \
					;; \
				linux) \
					os_arch="386 amd64 arm arm64 ppc64 ppc64le mips mipsle mips64 mips64le "; \
					;; \
				freebsd) \
					os_arch="386 amd64 arm arm64"; \
					;; \
				netbsd) \
					os_arch="386 amd64 arm arm64"; \
					;; \
				openbsd) \
					os_arch="386 amd64 arm arm64"; \
					;; \
				solaris) \
					os_arch="amd64"; \
					;; \
				windows) \
					os_arch="386 amd64"; \
					;; \
			esac; \
		fi; \
		for arch in $$os_arch; do \
			mkdir -p $(BIN_DIR)/$$os-$$arch; \
			echo "compiling for $$os/$$arch"; \
			env GOOS="$$os" GOARCH="$$arch" go build -o $(BIN_DIR)/$$os-$$arch/prometheus-html-exporter; \
		done ;\
	done

.PHONY: tarball
tarball:
	mkdir -p $(TARBALL_DIR)
	cd $(BIN_DIR); \
	for dir in $$(ls -d -1 */); do \
		echo "packaging $${dir%/}"; \
		mkdir -p $(TARBALL_DIR)/$$dir; \
		cp $(EXAMPLES_DIR)/* $$dir/* $(ROOT_DIR)/README.md $(ROOT_DIR)/LICENSE $(TARBALL_DIR)/$$dir; \
		cd $(TARBALL_DIR)/$$dir; \
		tar -czf ../prometheus-html-exporter-$${dir%/}.tar.gz *; \
		cd $(BIN_DIR); \
		rm -rf $(TARBALL_DIR)/$$dir; \
	done

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):latest .

.PHONY: release
release: build tarball docker-build
