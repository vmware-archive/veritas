#!/usr/bin/env bash

source ~/.bashisms/s3_upload.bash

set -e

echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .
echo "Uploading..."
upload_to_s3 veritas
rm veritas

echo "Compiling for OS X..."
go build .
mv veritas veritas-osx
upload_to_s3 veritas-osx
rm veritas-osx
