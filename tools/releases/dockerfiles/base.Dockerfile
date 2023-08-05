FROM golang:1.19 

WORKDIR /xds-control-plane

ADD go.mod go.sum /build-app/
RUN go mod download
COPY . .
RUN make service