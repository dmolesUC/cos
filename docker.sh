#!/bin/sh

VERSION=$(git describe --tags `git rev-list --tags --max-count=1`)
IMG_PREFIX="silaradost/cos"
IMAGE_TAG=$IMG_PREFIX":"$VERSION

docker build -t $IMAGE_TAG .
docker push $IMAGE_TAG
