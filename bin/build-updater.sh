#!/usr/bin/env sh

arch="${1}"
arch_tag="$(echo "${arch}" | tr '/' '-')"
commit="$(git rev-parse HEAD)"
image="ghcr.io/pipetail/bottlerocket-updater/updater:${arch_tag}-${commit}"
echo "building image: ${image}"

docker buildx build \
    --platform="${arch}" \
    --tag="${image}" \
    --file=images/updater/Dockerfile \
    --output type=image,push=true .
