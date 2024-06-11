#!/usr/bin/env bash

echo "Building..."

if [ -d "output" ]; then
    echo "Cleaning output..."
    rm -rf output
fi
mkdir output
cd /build/demo4/src || exit
go build -o /build/output/app
echo "OK!"
