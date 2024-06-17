#!/usr/bin/env bash

echo "Building..."

if [ -d "output" ]; then
    echo "Cleaning output..."
    rm -rf output
fi

mkdir -p /build/output/services
cd /build/demo5/src || exit
# build gateway

go build -o /build/output/app
echo "build gateway success"

#build service
pushd services || exit # enter service dir
for i in *; do
  name="$i"
  capName="${name^}"
  cd "$i" || exit
  go build -o "/build/output/services/$i/${capName}Service"
  echo "build $i success"
  cd ..
done
echo "OK!"
