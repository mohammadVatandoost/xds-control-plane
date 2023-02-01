.PHONY: help generate test lint fmt dependencies clean check coverage service race .remove_empty_dirs .pre-check-go

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))

service: $(SRCS)
	go build -o $@ -ldflags="$(LD_FLAGS)" ./


fmt: ## to run `go fmt` on all source code
	gofmt -s -w $(SRCS)
