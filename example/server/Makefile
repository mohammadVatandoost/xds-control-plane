.PHONY: help generate test lint fmt dependencies clean check coverage service race .remove_empty_dirs .pre-check-go

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))

DOCKER_REPOSITORY = "mvatandoost"
PROJECT_NAME = "xds-grpc-server-example"
TAG = "latest"
HELM_REPO_ADDRESS = "https://mohammadVatandoost.github.io/helm-chart/"
HELM_REPO_NAME = "myhelmrepo"
VERSION = "dev"
NAMESPACE = "test"

service: $(SRCS)
	go build -o $@ -ldflags="$(LD_FLAGS)" ./


fmt: ## to run `go fmt` on all source code
	gofmt -s -w $(SRCS)

build-docker:
	docker build . --tag $(DOCKER_REPOSITORY)/$(PROJECT_NAME):$(TAG)

helm-lint:
	helm lint deployments/helm/$(PROJECT_NAME)

helm-package:
	helm repo add $(HELM_REPO_NAME) $(HELM_REPO_ADDRESS)
	helm dependency update ./deployments/helm/$(PROJECT_NAME)
	helm package --app-version=$(VERSION) ./deployments/helm/$(PROJECT_NAME)

helm-deploy:
	helm -n $(NAMESPACE) upgrade -i $(PROJECT_NAME) -f ./deployments/helm/$(PROJECT_NAME)/values.yaml *.tgz 

