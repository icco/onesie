#! /bin/bash

TAG="$TRAVIS_BRANCH";
if [ -n "$TRAVIS_TAG" ]; then
  TAG="$TRAVIS_TAG";
fi

echo "BRANCH: $TRAVIS_BRANCH"
echo "TAG: $TRAVIS_TAG"
echo "Final tag: $TAG"

./packer build -var "account_file=onesie.json" -var "tag=$TAG" packer.json
