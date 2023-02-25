#!/bin/sh

CPU_ARCH="amd64"
if [[ $(uname -m) == *"arm"* ]]; then
  CPU_ARCH="aarch64"
fi

docker build --build-arg "CPU_ARCH=${CPU_ARCH}" -t gin-api .

echo "cleaning up old docker images..."
docker system prune -f
