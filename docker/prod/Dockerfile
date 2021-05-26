FROM golang:1.15-alpine as build-stage

RUN apk --no-cache add \
    g++ \
    git \
    make \
    bash \
    upx

ARG VERSION
ENV VERSION=${VERSION}

WORKDIR /src
COPY . .
RUN ./scripts/build/build.sh
RUN upx /src/bin/k8s-webhook-pull-policy

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates
COPY --from=build-stage /src/bin/k8s-webhook-pull-policy /usr/local/bin/k8s-webhook-pull-policy
ENTRYPOINT ["/usr/local/bin/k8s-webhook-pull-policy"]