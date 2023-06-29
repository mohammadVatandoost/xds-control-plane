.PHONY: build help generate test lint fmt dependencies clean check coverage service race .remove_empty_dirs .pre-check-go

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))

DOCKER_REPOSITORY = "mvatandoost"
PROJECT_NAME = "xds-control-plane"
TAG = "latest"
HELM_REPO_ADDRESS = "https://mohammadVatandoost.github.io/helm-chart/"
HELM_REPO_NAME = "myhelmrepo"
VERSION = "dev"
NAMESPACE = "test"

build: $(SRCS)
	go build -o ./build/$(PROJECT_NAME) -ldflags="$(LD_FLAGS)" ./cmd/...


fmt: ## to run `go fmt` on all source code
	gofmt -s -w $(SRCS)

build-image:
	docker build . -f build/Dockerfile --tag $(DOCKER_REPOSITORY)/$(PROJECT_NAME):$(TAG)

kind-load:
	kind load docker-image $(DOCKER_REPOSITORY)/$(PROJECT_NAME):$(TAG)	

helm-lint:
	helm lint deployments/helm/$(PROJECT_NAME)

helm-package:
	helm repo add $(HELM_REPO_NAME) $(HELM_REPO_ADDRESS)
	helm dependency update ./deployments/helm/$(PROJECT_NAME)
	helm package --app-version=$(VERSION) ./deployments/helm/$(PROJECT_NAME)

helm-deploy:
	helm -n $(NAMESPACE) upgrade -i $(PROJECT_NAME) -f ./deployments/helm/$(PROJECT_NAME)/values.yaml *.tgz 

helm-ci-cd:
	helm lint deployments/helm/$(PROJECT_NAME)
	helm repo add $(HELM_REPO_NAME) $(HELM_REPO_ADDRESS)
	helm dependency update ./deployments/helm/$(PROJECT_NAME)
	helm package --app-version=$(VERSION) ./deployments/helm/$(PROJECT_NAME)
	helm -n $(NAMESPACE) upgrade -i $(PROJECT_NAME) -f ./deployments/helm/$(PROJECT_NAME)/values.yaml *.tgz
