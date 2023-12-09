#!/bin/bash

set -e

TOKEN=$(yc iam create-token)
docker login -u iam -p ${TOKEN} cr.yandex

set -ex

REGISTRY="cr.yandex/crp7b3092gvddqsd8k3u"
SERVICE="$1"


docker build --platform linux/amd64 --rm -t ${REGISTRY}/${SERVICE}:latest -f ./infra/k8s/${SERVICE}/${SERVICE}.Dockerfile .

docker push ${REGISTRY}/${SERVICE}:latest
