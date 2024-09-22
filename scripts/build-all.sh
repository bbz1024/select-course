#!/usr/bin/env bash

echo "Building..."

if [ -d "output" ]; then
    echo "Cleaning output..."
    rm -rf output
fi
# 压缩
compress=${compress:-0}  # 使用默认值0，如果未定义或为空

mkdir -p /build/output/services
cd /build/demo7/src || exit
# build gateway

go build -o /build/output/app
if [ "$compress" == 1 ]; then
    echo "compress..."
    /build/pak/upx -9 /build/output/app
fi
echo "build gateway success"

#build service


pushd services || exit # enter service dir
for i in *; do
  name="$i"
  capName="${name^}"
  cd "$i" || exit
  go build -o "/build/output/services/$i/${capName}Service"
  if [ "$compress" == 1 ]; then
        echo "compress..."
        /build/pak/upx -9 "/build/output/services/$i/${capName}Service"
  fi
  echo "build $i success"
  cd ..
done
echo "OK!"
