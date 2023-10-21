ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH as base_image
ARG ARCH

FROM debian:stable-20220509 as final
ARG ARCH


COPY --from=base_image /control-plane/build/artifacts-linux-${ARCH}/cp-example-server/cp-example-server /usr/bin

WORKDIR /usr/bin
ENTRYPOINT ["./cp-example-server"]