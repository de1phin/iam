#!/bin/bash

id=$(yc lockbox secret list --folder-id=b1gqe3skkuiko3bv671e --format=json | jq -r '.[] | select(.name == "alexzhaba-sa") | .id')
token=$(yc iam create-token)
secret=$(curl -s -H "Authorization: Bearer $token" https://payload.lockbox.api.cloud.yandex.net/lockbox/v1/secrets/$id/payload)
export AWS_ACCESS_KEY_ID=$(echo $secret | jq -r '.entries[] | select(.key == "key-id") | .textValue')
export AWS_SECRET_ACCESS_KEY=$(echo $secret | jq -r '.entries[] | select(.key == "secret-key") | .textValue')
