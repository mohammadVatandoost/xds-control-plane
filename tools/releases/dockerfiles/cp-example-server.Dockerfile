ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH
ARG ARCH

WORKDIR /control-plane

COPY ./build/artifacts-linux-${ARCH}/cp-example-server/cp-example-server /usr/bin

ENTRYPOINT ["/usr/bin/cp-example-server"]