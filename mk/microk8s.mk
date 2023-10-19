
.PHONY: microk8s/deploy/control-plane
microk8s/deploy/control-plane: microk8s/load
	microk8s helm upgrade --install --namespace $(CONTROL_PLANE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-control-plane ./deployments/helm/xds-control-plane
	microk8s helm upgrade --install --namespace $(EXAMPLE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-grpc-client-example ./example/client/deployments/helm/xds-grpc-client-example
	microk8s helm upgrade --install --namespace $(EXAMPLE_NAMESPACE) --create-namespace \
                --set global.image.registry="$(DOCKER_REGISTRY)" \
                --set global.image.tag="$(BUILD_INFO_VERSION)" \
				xds-grpc-server-example ./example/server/deployments/helm/xds-grpc-server-example	

.PHONY: microk8s/load/images
microk8s/load/images:
	for image in ${CONTROL_PLANE_IMAGES}; do docker save $$image > "$$image".tar; done
	for image in ${CONTROL_PLANE_IMAGES}; do microk8s ctr image import "$$image".tar; done

.PHONY: microk8s/load
microk8s/load: images docker/tag microk8s/load/images