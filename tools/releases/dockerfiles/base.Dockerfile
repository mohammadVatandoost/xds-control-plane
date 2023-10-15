FROM golang:1.21 

WORKDIR /control-plane

ADD go.mod go.sum /control-plane/
RUN go mod download
COPY . .
RUN make build/release