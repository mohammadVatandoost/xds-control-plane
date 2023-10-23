EXAMPLE_NAMESPACE ?= control-plane-example
CONTROL_PLANE_NAMESPACE ?= control-plane
KIND_CLUSTER_NAME ?= control-plane

# The e2e tests depend on Kind kubeconfigs being in this directory,
# so this is location should not be changed by developers.
KIND_KUBECONFIG_DIR := $(HOME)/.kube

# This is the name of the current config file to use.
KIND_KUBECONFIG := $(KIND_KUBECONFIG_DIR)/kind-$(KIND_CLUSTER_NAME)-config

# Ensure Kubernetes tooling only gets the config we explicity specify.
unexport KUBECONFIG

METRICS_SERVER_VERSION := 0.4.1

ifdef IPV6
KIND_CONFIG ?= $(TOP)/test/kind/cluster-ipv6.yaml
else
KIND_CONFIG ?= $(TOP)/test/kind/cluster.yaml
endif


CI_KUBERNETES_VERSION ?= v1.22.9@sha256:8135260b959dfe320206eb36b3aeda9cffcb262f4b44cda6b33f7bb73f453105

CONTROL_PLANE_MODE ?= standalone
CONTROL_PLANE_NAMESPACE ?= control-plane-system


.PHONY: kind/start
kind/start: ${KUBECONFIG_DIR}
	@$(KIND) get clusters | grep $(KIND_CLUSTER_NAME) >/dev/null 2>&1 && echo "Kind cluster already running." && exit 0 || \
		($(KIND) create cluster \
			--name "$(KIND_CLUSTER_NAME)" \
			--kubeconfig $(KIND_KUBECONFIG) \
			--quiet --wait 120s && \
		KUBECONFIG=$(KIND_KUBECONFIG) $(KUBECTL) scale deployment --replicas 1 coredns --namespace kube-system && \
		$(MAKE) kind/wait)
	@echo
	@echo '>>> You need to manually run the following command in your shell: >>>'
	@echo
	@echo export KUBECONFIG="$(KIND_KUBECONFIG)"
	@echo
	@echo '<<< ------------------------------------------------------------- <<<'
	@echo


.PHONY: kind/wait
kind/wait:
	until \
		KUBECONFIG=$(KIND_KUBECONFIG) $(KUBECTL) wait -n kube-system --timeout=5s --for condition=Ready --all pods ; \
	do echo "Waiting for the cluster to come up" && sleep 1; done

.PHONY: kind/stop
kind/stop:
	@$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)
	@rm -f $(KUBECONFIG_DIR)/$(KIND_KUBECONFIG)

.PHONY: kind/stop/all
kind/stop/all:
	@$(KIND) delete clusters --all
	@rm -f $(KUBECONFIG_DIR)/kind-control-plane-*

.PHONY: kind/load/images
kind/load/images:
	for image in ${CONTROL_PLANE_IMAGES}; do $(KIND) load docker-image $$image --name=$(KIND_CLUSTER_NAME); done

.PHONY: kind/load
kind/load: images docker/tag kind/load/images


.PHONY: kind/deploy/control-plane
kind/deploy/control-plane: kind/load
	KUBECONFIG=$(KIND_KUBECONFIG) helm upgrade --install --namespace $(CONTROL_PLANE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-control-plane ./deployments/helm/xds-control-plane
	KUBECONFIG=$(KIND_KUBECONFIG) helm upgrade --install --namespace $(EXAMPLE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-grpc-client-example ./example/client/deployments/helm/xds-grpc-client-example
	KUBECONFIG=$(KIND_KUBECONFIG) helm upgrade --install --namespace $(EXAMPLE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-grpc-server-example ./example/server/deployments/helm/xds-grpc-server-example						


.PHONY: kind/delete/control-plane
kind/delete/control-plane:
	KUBECONFIG=$(KIND_KUBECONFIG) helm uninstall --namespace $(CONTROL_PLANE_NAMESPACE) xds-control-plane | true
	KUBECONFIG=$(KIND_KUBECONFIG) helm uninstall  --namespace $(EXAMPLE_NAMESPACE) xds-grpc-client-example | true
	KUBECONFIG=$(KIND_KUBECONFIG) helm uninstall --namespace $(EXAMPLE_NAMESPACE) xds-grpc-server-example | true
	KUBECONFIG=$(KIND_KUBECONFIG) $(KUBECTL) delete namespace $(CONTROL_PLANE_NAMESPACE) | true
	KUBECONFIG=$(KIND_KUBECONFIG) $(KUBECTL) delete namespace $(EXAMPLE_NAMESPACE) | true	


