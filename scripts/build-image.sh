#!/usr/bin/env bash

set -e

[ -z "$VERSION" ] && echo "VERSION env var is required." && exit 1;
[ -z "$IMAGE" ] && echo "IMAGE env var is required." && exit 1;

# By default use amd64 architecture.
ARCH=${ARCH:-amd64}
DOCKER_FILE=${DOCKER_FILE:-Dockerfile.ci}

IMAGE_TAG_ARCH="${IMAGE}:${VERSION}-${ARCH}"

# Build image.
echo "Building image ${IMAGE_TAG_ARCH}..."
docker build \
    --build-arg VERSION="${VERSION}" \
    --build-arg ARCH="${ARCH}" \
    -f ${DOCKER_FILE} \
    -t "${IMAGE_TAG_ARCH}" \
    -t "${IMAGE}:latest" .

# Push to registry
if [ "$IMAGE" != "juggler-dev" ]
then
  echo "Pushing latest tag and ${IMAGE_TAG_ARCH}..."
  docker push ${IMAGE_TAG_ARCH}
  docker push ${IMAGE}:latest
else
  docker tag ${IMAGE_TAG_ARCH} $IMAGE:latest
fi