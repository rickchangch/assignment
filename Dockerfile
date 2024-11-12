FROM golang:1.22.8 as builder
WORKDIR /workspace
RUN apt-get update && \
    apt-get install -y make curl && \
    rm -rf /var/lib/apt/lists/*
ADD . /workspace
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make build

FROM alpine:3.14
WORKDIR /workspace
RUN apk upgrade && \
    apk add --no-cache ca-certificates=20230506-r0 && \
   apk --no-cache add shadow
COPY --from=builder /workspace/build ./build
COPY --from=builder /workspace/setup/config ./setup/config

ENTRYPOINT ["/workspace/build/server"]