ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH
ARG ARCH

WORKDIR /control-plane

COPY ./build/artifacts-linux-${ARCH}/cp-example-client/cp-example-client /usr/bin
COPY ./example/client/xds_bootstrap.json /usr/bin
COPY ./example/client/xds_bootstrap_local.json /usr/bin
WORKDIR /usr/bin
ENTRYPOINT ["./cp-example-client"]