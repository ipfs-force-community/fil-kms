SHELL=/usr/bin/env bash

all: build
.PHONY: all


unexport GOFLAGS

GOCC?=go

MODULES:=


## FFI

FFI_PATH:=extern/filecoin-ffi/
FFI_DEPS:=.install-filcrypto
FFI_DEPS:=$(addprefix $(FFI_PATH),$(FFI_DEPS))

$(FFI_DEPS): build/.filecoin-install ;

build/.filecoin-install: $(FFI_PATH)
	$(MAKE) -C $(FFI_PATH) $(FFI_DEPS:$(FFI_PATH)%=%)

MODULES+=$(FFI_PATH)
BUILD_DEPS+=build/.filecoin-install
CLEAN+=build/.filecoin-install

ffi-version-check:
	@[[ "$$(awk '/const Version/{print $$5}' extern/filecoin-ffi/version.go)" -eq 3 ]] || (echo "FFI version mismatch, update submodules"; exit 1)
BUILD_DEPS+=ffi-version-check

.PHONY: ffi-version-check



build: fil-kms
.PHONY: build

clean:
	rm -rf fil-kms
	go clean
	-$(MAKE) -C ./extern/filecoin-ffi clean



fil-kms:  $(BUILD_DEPS)
	rm -f fil-kms
	go build -o fil-kms .
.PHONY: fil-kms