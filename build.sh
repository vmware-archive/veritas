#!/usr/bin/env bash

source ~/.bash_it/bash_it.sh

set -e

echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .
echo "Uploading..."
upload_to_s3 veritas
rm veritas
