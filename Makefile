
# os detection
ifeq (${OS},Windows_NT)
MKDIR = mkdir $(subst /,\,$(1)) > nul 2>&1 || (exit 0)
WHICH := where
DEVNULL := NUL
else
MKDIR = mkdir -p $(1)
WHICH := which
DEVNULL := /dev/null
endif

ifneq ($(shell $(WHICH) cygpath 2>$(DEVNULL)),)
OSCURDIR := $(shell cygpath -w ${CURDIR})
else
ifneq ($(shell $(WHICH) wslpath 2>$(DEVNULL)),)
OSCURDIR := $(shell wslpath -w ${CURDIR})
else
OSCURDIR := ${CURDIR}
endif
endif

# detect go
ifneq ($(shell $(WHICH) go 2>$(DEVNULL)),)
GO := go
else 
ifneq ($(shell $(WHICH) go.exe 2>$(DEVNULL)),)
GO := go.exe
else
$(error "go is not in your system PATH")
endif
endif

# Cpmpute the list of packages to include in coverage analysis, excluding any package that contains "pb" in its path
COVERPKGS := $(shell $(GO) list ./... | grep -E -v "pb" | paste -sd ",")

.PHONY: all
all: build test

.PHONY: build
build:
	@echo "[build]"
	@$(GO) build ./...

.PHONY: test
test:
	@echo "[test]"
	@$(GO) test ./...

.PHONY: lint
lint:
	@echo "[lint]"
	@docker run --rm -v $(CURDIR):/pwd -w /pwd golangci/golangci-lint:v2.9.0-alpine golangci-lint run --timeout=5m

.PHONY: coverage
coverage:
	@echo "[coverage]"
	@rm -rf coverage/percent
	@${MKDIR} coverage/percent/data
	@echo "[test]"
	@$(GO) test -cover ./... -coverpkg=${COVERPKGS} -args -test.gocoverdir="${OSCURDIR}/coverage/percent/data" 1>$(DEVNULL)
	@echo "[cover]"
	@$(GO) tool covdata percent -i=./coverage/percent/data

.PHONY: coverage-html
coverage-html:
	@echo "[coverage-html]"
	@rm -rf coverage/html
	@${MKDIR} coverage/html/data
	@echo "[test]"
	@$(GO) test -cover ./... -coverpkg=${COVERPKGS} -args -test.gocoverdir="${OSCURDIR}/coverage/html/data" 1>$(DEVNULL)
	@echo "[cover]"
	@$(GO) tool covdata textfmt -i=./coverage/html/data -o coverage/html/profile 1>$(DEVNULL)
	@$(GO) tool cover -html coverage/html/profile -o coverage/html/coverage.html

