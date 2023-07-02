FROM golang:1.19 as base_image

WORKDIR /build-app

ADD go.mod go.sum /build-app/
RUN go mod download
COPY . .
RUN make service



FROM debian:stable-20220509 as final

WORKDIR /app

COPY --from=base_image /build-app/service .
COPY --from=base_image /build-app/xds_bootstrap.json .
COPY --from=base_image /build-app/tls .


ENTRYPOINT ["./service"]
