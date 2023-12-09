#!/bin/bash

# for now use same preset password for all databases
if [[ -z "${PSQL_PASSWORD}" ]]; then
    exit 1
fi

jq -n --arg password "$PSQL_PASSWORD" \
    '{"password": $password}'
