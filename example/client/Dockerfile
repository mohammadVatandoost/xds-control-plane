FROM golang:1.19 as base_image

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build-app

ADD go.mod go.sum /build-app/
RUN go mod download
COPY . .
RUN make service



FROM debian:stable-20220509 as final

WORKDIR /app

COPY --from=base_image /build-app/service .

ENTRYPOINT ["./service"]