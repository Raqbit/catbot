#!/usr/bin/env bash

# TODO: update for production?
source ./creds.sh

docker run \
  -it \
  -u "$UID:$GID" \
  --rm \
  -v "$PWD/migrations:/migrations" \
  --network host \
  migrate/migrate \
  -path=/migrations/ \
  -database "$SOURCE" \
  down
