#! /bin/bash

TAG="$TRAVIS_BRANCH";
if [ -z "$TRAVIS_TAG" ]; then
  TAG="$TRAVIS_TAG";
fi

./packer build -var "account_file=onesie.json" -var "tag=$TAG" packer.json
