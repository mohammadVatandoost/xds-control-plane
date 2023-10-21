ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH as base_image
ARG ARCH

FROM debian:stable-20220509 as final
ARG ARCH


COPY --from=base_image /control-plane/build/artifacts-linux-${ARCH}/control-plane/control-plane /usr/bin

WORKDIR /usr/bin

ENTRYPOINT ["./control-plane"]