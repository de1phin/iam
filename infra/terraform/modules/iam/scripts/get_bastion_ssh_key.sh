#!/bin/bash

ssh_key=$(ssh -i ${BASTION_SSH_KEY} ${BASTION_HOST} cat ${BASTION_ROOT_SSH_KEY_FILE})

jq -n --arg key "${ssh_key}" \
    '{"key": $key}'
