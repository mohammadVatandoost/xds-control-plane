ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH as base_image
ARG ARCH

FROM debian:stable-20220509 as final
ARG ARCH


COPY --from=base_image /control-plane/build/artifacts-linux-${ARCH}/cp-example-client/cp-example-client /usr/bin
COPY --from=base_image /control-plane/example/client/xds_bootstrap.json /usr/bin
COPY --from=base_image /control-plane/example/client/xds_bootstrap_local.json /usr/bin

WORKDIR /usr/bin
ENTRYPOINT ["./cp-example-client"]