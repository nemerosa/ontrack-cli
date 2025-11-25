#!/usr/bin/env bash

TAG=$1
if [[ "$TAG" == "" ]]; then
  echo "Tag parameter is required"
  exit 1
fi

git tag "$TAG"
git push origin "$TAG"
