DEFAULT_OUT=bin/kacpersh
ifdef GOOS
DEFAULT_OUT:=$(DEFAULT_OUT).$(GOOS)
endif
ifdef GOARCH
DEFAULT_OUT:=$(DEFAULT_OUT).$(GOARCH)
endif
ifeq ($(GOOS),windows)
DEFAULT_OUT:=$(DEFAULT_OUT).exe
endif
OUT?=$(DEFAULT_OUT)

.PHONY: build
build:
	$(info $(OUTT))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUT) 

.PHONY: build-linux
build-linux:
	$(MAKE) build GOOS=linux GOARCH=amd64 OUT=bin/kacpersh

.PHONY: lint
lint:
	golangci-lint run