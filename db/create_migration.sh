#!/usr/bin/env bash

if [ "$#" -ne 1 ]; then
    echo "Please specify the name of the migration to create"
    exit 2
fi

USER=$(id -u "$USER")
GROUP=$(id -g "$USER")

docker run \
  -u "$USER:$GROUP" \
  --rm \
  -v "$PWD/migrations:/migrations" \
  --network host \
  migrate/migrate \
  create \
  -ext sql \
  -dir /migrations \
  -seq "$1"
