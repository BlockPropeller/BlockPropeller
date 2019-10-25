#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Please specify image tag"
  exit 1
fi

TAG="$1"

docker build -t "gcr.io/tenderly-project/blockpropeller-api:${TAG}" \
  -t "gcr.io/tenderly-project/blockpropeller-api:latest" .

docker push "gcr.io/tenderly-project/blockpropeller-api:${TAG}"
docker push "gcr.io/tenderly-project/blockpropeller-api:latest"
