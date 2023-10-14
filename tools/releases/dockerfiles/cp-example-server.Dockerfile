ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH
ARG ARCH

WORKDIR /control-plane

COPY ./build/artifacts-linux-${ARCH}/example-server/example-server /usr/bin

ENTRYPOINT ["/usr/bin/example-server"]