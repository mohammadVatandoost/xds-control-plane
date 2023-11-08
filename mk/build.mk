build_info_fields = \
	version=$(BUILD_INFO_VERSION) \
	gitTag=$(GIT_TAG) \
	gitCommit=$(GIT_COMMIT) \
	buildDate=$(BUILD_DATE) 

build_info_ld_flags := $(foreach entry,$(build_info_fields), -X github.com/mohammadVatandoost/xds-conrol-plane/pkg/version.$(entry))

LD_FLAGS := -ldflags="-s -w $(build_info_ld_flags) $(EXTRA_LD_FLAGS)"
EXTRA_GOENV ?=
GOENV=CGO_ENABLED=0 $(EXTRA_GOENV)
GOFLAGS := -trimpath $(EXTRA_GOFLAGS)

TOP := $(shell pwd)
BUILD_DIR ?= $(TOP)/build


BUILD_RELEASE_BINARIES := control-plane cp-example-client cp-example-server 

SUPPORTED_GOARCHES ?= amd64 arm64

SUPPORTED_GOOSES ?= linux darwin

ENABLED_GOARCHES ?= $(GOARCH)
ENABLED_GOOSES ?= $(GOOS)

IGNORED_ARCH_OS ?=
ENABLED_ARCH_OS = $(filter-out $(IGNORED_ARCH_OS), $(foreach os,$(ENABLED_GOOSES),$(foreach arch,$(ENABLED_GOARCHES),$(os)-$(arch))))

.PHONY: build/info
build/info: 
	@echo build-info: $(build_info_fields)
	@echo tools-dir: $(CI_TOOLS_DIR)
	@echo arch: supported=$(SUPPORTED_GOARCHES), enabled=$(ENABLED_GOARCHES)
	@echo os: supported=$(SUPPORTED_GOOSES), enabled=$(ENABLED_GOOSES)
	@echo ignored=$(IGNORED_ARCH_OS)
	@echo enabled arch-os=$(ENABLED_ARCH_OS)

.PHONY: build
build: build/release

.PHONY: build/release
build/release: $(addprefix build/,$(BUILD_RELEASE_BINARIES)) 

define LOCAL_BUILD_TARGET
build/$(1): $$(patsubst %,build/artifacts-%/$(1),$$(ENABLED_ARCH_OS))
endef
$(foreach target,$(BUILD_RELEASE_BINARIES),$(eval $(call LOCAL_BUILD_TARGET,$(target))))

Build_Go_Application = GOOS=$(1) GOARCH=$(2) $$(GOENV) go build -v $$(GOFLAGS) $$(LD_FLAGS) -o $$@/$$(notdir $$@)


define BUILD_TARGET
.PHONY: build/artifacts-$(1)-$(2)/control-plane
build/artifacts-$(1)-$(2)/control-plane:
	$(Build_Go_Application) ./cmd

.PHONY: build/artifacts-$(1)-$(2)/cp-example-client
build/artifacts-$(1)-$(2)/cp-example-client:
	$(Build_Go_Application) ./example/client

.PHONY: build/artifacts-$(1)-$(2)/cp-example-server-1
build/artifacts-$(1)-$(2)/cp-example-server-1:
	$(Build_Go_Application) ./example/server1

.PHONY: build/artifacts-$(1)-$(2)/cp-example-server-2
build/artifacts-$(1)-$(2)/cp-example-server-2:
	$(Build_Go_Application) ./example/server2	

endef
$(foreach goos,$(SUPPORTED_GOOSES),$(foreach goarch,$(SUPPORTED_GOARCHES),$(eval $(call BUILD_TARGET,$(goos),$(goarch)))))

.PHONY: clean/build
clean/build: # clean/ebpf ## Dev: Remove build/ dir
	rm -rf "$(BUILD_DIR)"
