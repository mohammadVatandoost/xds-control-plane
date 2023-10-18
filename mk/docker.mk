BUILD_DOCKER_IMAGES_DIR ?= $(BUILD_DIR)/docker-images-${GOARCH}
CONTROL_PLANE_VERSION ?= master

DOCKER_SERVER ?= docker.io
# DOCKER_REGISTRY ?= $(DOCKER_SERVER)/mvatandoost
DOCKER_REGISTRY ?= mvatandoost
DOCKER_USERNAME ?=
DOCKER_API_KEY ?=

define build_image
$(addsuffix :$(BUILD_INFO_VERSION)$(if $(2),-$(2)),$(addprefix $(DOCKER_REGISTRY)/,$(1)))
endef

IMAGES_RELEASE += control-plane cp-example-server cp-example-client
CONTROL_PLANE_IMAGES = $(call build_image,$(IMAGES_RELEASE))

.PHONY: images/show
images/show: ## output all images that are built with the current configuration
	@echo $(CONTROL_PLANE_IMAGES)

# Always use Docker BuildKit, see
# https://docs.docker.com/develop/develop-images/build_enhancements/
export DOCKER_BUILDKIT := 1

# add targets to build images for each arch
# $(1) - GOOS to build for
define IMAGE_TARGETS_BY_ARCH
.PHONY: image/base/$(1)
image/base/$(1): ## Dev: Rebuild `control-plane-base` Docker image
	docker build -t control-plane/base-nossl-debian11:no-push-$(1) --build-arg ARCH=$(1) --platform=linux/$(1) -f $(TOOLS_DIR)/releases/dockerfiles/base.Dockerfile .

.PHONY: image/control-plane/$(1)
image/control-plane/$(1): image/base/$(1) build/artifacts-linux-$(1)/control-plane ## Dev: Rebuild `control-plane` Docker image
	docker build -t $$(call build_image,control-plane,$(1)) --build-arg ARCH=$(1) --platform=linux/$(1) -f $(TOOLS_DIR)/releases/dockerfiles/control-plane.Dockerfile .

.PHONY: image/cp-example-client/$(1)
image/cp-example-client/$(1): image/base/$(1) build/artifacts-linux-$(1)/cp-example-client ## Dev: Rebuild `cp-example-client` Docker image
	docker build -t $$(call build_image,cp-example-client,$(1)) --build-arg ARCH=$(1) --platform=linux/$(1) -f $(TOOLS_DIR)/releases/dockerfiles/cp-example-client.Dockerfile .

.PHONY: image/cp-example-server/$(1)
image/cp-example-server/$(1): image/base/$(1) build/artifacts-linux-$(1)/cp-example-server ## Dev: Rebuild `cp-example-server` Docker image
	docker build -t $$(call build_image,cp-example-server,$(1)) --build-arg ARCH=$(1) --platform=linux/$(1) -f $(TOOLS_DIR)/releases/dockerfiles/cp-example-server.Dockerfile .

endef
$(foreach goarch,$(SUPPORTED_GOARCHES),$(eval $(call IMAGE_TARGETS_BY_ARCH,$(goarch))))

# add targets to generate docker/{save,load,tag,push} for each supported ARCH
# add targets to build images for each arch
# $(1) - GOOS to build for
# $(2) - GOARCH to build for
define DOCKER_TARGETS_BY_ARCH
.PHONY: docker/$(1)/$(2)/save
docker/$(1)/$(2)/save:
	@mkdir -p build/docker
	docker save --output build/docker/$(1)-$(2).tar $$(call build_image,$(1),$(2))

.PHONY: docker/$(1)/$(2)/load
docker/$(1)/$(2)/load:
	@docker load --quiet --input build/docker/$(1)-$(2).tar

# we only tag the image that has the same arch than the HOST (tag is meant to use the image just after so having the same arch makes sense)
.PHONY: docker/$(1)/$(2)/tag
docker/$(1)/$(2)/tag:
	$$(if $$(findstring $(GOARCH),$(2)),docker tag $$(call build_image,$(1),$(2)) $$(call build_image,$(1)),# Not tagging $(1) as $(2) is not host arch)

.PHONY: docker/$(1)/$(2)/push
docker/$(1)/$(2)/push:
	$$(call GATE_PUSH,docker push $$(call build_image,$(1),$(2)))
endef
$(foreach goarch, $(SUPPORTED_GOARCHES),$(foreach image, $(IMAGES_RELEASE) $(IMAGES_TEST),$(eval $(call DOCKER_TARGETS_BY_ARCH,$(image),$(goarch)))))

# create and push a manifest for each
docker/%/manifest:
	$(call GATE_PUSH,docker manifest create $(call build_image,$*) $(patsubst %,--amend $(call build_image,$*,%),$(ENABLED_GOARCHES)))
	$(call GATE_PUSH,docker manifest push $(call build_image,$*))

# add targets like `docker/save` with dependencies all `ENABLED_GOARCHES`
ALL_RELEASE_WITH_ARCH=$(foreach arch,$(ENABLED_GOARCHES),$(patsubst %,%/$(arch),$(IMAGES_RELEASE)))

test:
	echo $(patsubst %,docker/%/tag,$(ALL_RELEASE_WITH_ARCH))

ALL_TEST_WITH_ARCH=$(foreach arch,$(ENABLED_GOARCHES),$(patsubst %,%/$(arch),$(IMAGES_TEST)))
.PHONY: docker/save
docker/save: $(patsubst %,docker/%/save,$(ALL_RELEASE_WITH_ARCH) $(ALL_TEST_WITH_ARCH))
.PHONY: docker/load
docker/load: $(patsubst %,docker/%/load,$(ALL_RELEASE_WITH_ARCH) $(ALL_TEST_WITH_ARCH))
.PHONY: docker/tag
docker/tag: docker/tag/test docker/tag/release ## Tag local arch containers with the version with the arch (this is mostly to use non multi-arch images as if they were released images in e2e tests)
.PHONY: docker/tag/release
docker/tag/release: $(patsubst %,docker/%/tag,$(ALL_RELEASE_WITH_ARCH))
.PHONY: docker/tag/test
docker/tag/test: $(patsubst %,docker/%/tag,$(ALL_TEST_WITH_ARCH))
.PHONY: docker/push
docker/push: $(patsubst %,docker/%/push,$(ALL_RELEASE_WITH_ARCH)) ## Publish all docker images with arch specific tags
.PHONY: docker/manifest
docker/manifest: $(patsubst %,docker/%/manifest,$(IMAGES_RELEASE)) ## Publish all manifests (images need to be pushed already
.PHONY: images
images: images/release  ## Dev: Rebuild release and test Docker images
.PHONY: images/release
images/release: $(addprefix image/,$(ALL_RELEASE_WITH_ARCH)) ## Dev: Rebuild release Docker images


.PHONY: docker/purge
docker/purge: ## Dev: Remove all Docker containers, images, networks and volumes
	for c in `docker ps -q`; do docker kill $$c; done
	docker system prune --all --volumes --force

.PHONY: docker/login
docker/login:
	$(call GATE_PUSH,docker login -u $(DOCKER_USERNAME) -p $(DOCKER_API_KEY) $(DOCKER_SERVER))

.PHONY: docker/logout
docker/logout:
	$(call GATE_PUSH,docker logout $(DOCKER_SERVER))
