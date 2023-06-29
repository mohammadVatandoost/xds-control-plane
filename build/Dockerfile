FROM golang:1.19 as base_image

WORKDIR /build-app

ADD go.mod go.sum /build-app/
RUN go mod download
COPY . .
RUN make build




FROM debian:stable-20220509 as final

WORKDIR /app

COPY --from=base_image /build-app/build/xds-control-plane .

ENTRYPOINT ["./xds-control-plane", "serve"]