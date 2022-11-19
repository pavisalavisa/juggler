# Set also `ARCH` ARG here so we can use it on all the `FROM`s.
ARG ARCH

FROM golang:1.19-alpine as build-go

RUN apk --no-cache add \
    g++ \
    git \
    make \
    curl \
    bash

# Required by the built script for setting verion and cross-compiling.
ARG VERSION
ENV VERSION=${VERSION}
ARG ARCH
ENV GOARCH=${ARCH}

# Copy go.mod and go.sum to use docker cache
COPY go.mod go.sum /code/
WORKDIR /code/
RUN go mod download

# Copy the source files
COPY . /code/
RUN /code/scripts/build.sh juggler

FROM alpine:3
COPY --from=build-go /code/bin/juggler /juggler
ENTRYPOINT ["/juggler"]