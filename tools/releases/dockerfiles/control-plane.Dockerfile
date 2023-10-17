ARG ARCH
FROM control-plane/base-nossl-debian11:no-push-$ARCH
ARG ARCH

WORKDIR /control-plane

COPY ./build/artifacts-linux-${ARCH}/control-plane/control-plane /usr/bin

ENTRYPOINT ["/usr/bin/control-plane"]