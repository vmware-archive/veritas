#!/usr/bin/env bash

if [ -z "$GITHUB_API_TOKEN" ]; then
    echo "Please set GITHUB_API_TOKEN"
    echo "To create a token, visit: https://github.com/settings/tokens"
    exit 1
fi  

set -e

RELEASE_ID=`curl -s -H "Authorization: token $GITHUB_API_TOKEN" https://api.github.com/repos/pivotal-cf-experimental/veritas/releases/latest | jq ".id"`
LINUX_ASSET_ID=`curl -s -H "Authorization: token $GITHUB_API_TOKEN" https://api.github.com/repos/pivotal-cf-experimental/veritas/releases/latest | jq ".assets[0].id"`
OSX_ASSET_ID=`curl -s -H "Authorization: token $GITHUB_API_TOKEN" https://api.github.com/repos/pivotal-cf-experimental/veritas/releases/latest | jq ".assets[1].id"`

echo "Deleting existing Linux Asset"
curl -H "Authorization: token $GITHUB_API_TOKEN" \
     -X DELETE \
     "https://api.github.com/repos/pivotal-cf-experimental/veritas/releases/assets/$LINUX_ASSET_ID"
echo "Deleting existing OS X Asset"
curl -H "Authorization: token $GITHUB_API_TOKEN" \
     -X DELETE \
     "https://api.github.com/repos/pivotal-cf-experimental/veritas/releases/assets/$OSX_ASSET_ID"


echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .
echo "Uploading..."
curl -H "Authorization: token $GITHUB_API_TOKEN" \
     -H "Content-Type: application/octet-stream" \
     --data-binary @veritas \
     "https://uploads.github.com/repos/pivotal-cf-experimental/veritas/releases/$RELEASE_ID/assets?name=veritas&label=Veritas%20%28Linux%29"
rm veritas

echo "Compiling for OS X..."
go build .
echo "Uploading..."
curl -H "Authorization: token $GITHUB_API_TOKEN" \
     -H "Content-Type: application/octet-stream" \
     --data-binary @veritas \
     "https://uploads.github.com/repos/pivotal-cf-experimental/veritas/releases/$RELEASE_ID/assets?name=veritas-osx&label=Veritas%20%28OS%20X%29"
rm veritas
