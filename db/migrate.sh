#!/usr/bin/env bash

function migrate() {
    docker run \
      -u "$UID:$GID" \
      --rm \
      -v "$PWD/migrations:/migrations" \
      --network host \
      migrate/migrate \
      -path=/migrations/ \
      -database "$1" \
      up
}

source ./creds.sh

if [[ $MIGRATION_MODE == "production" ]]
then
  read -r -p "Are you sure you want to migrate production? [y/N] " response
  if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
  then
      if [[ $MIGRATION_PRODUCTION_URL ]]
      then
        SOURCE=$MIGRATION_PRODUCTION_URL
      else
        echo "Please set the MIGRATION_PRODUCTION_URL environment variable to the appropriate url."
        exit 1
      fi
  else
      echo "Not migrating production."
      exit 1
  fi
fi

migrate "$SOURCE"


