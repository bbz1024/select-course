#!/bin/bash

if [ ! "$(basename "$PWD")" ] = "select-course"; then
  echo "please run this script in select-course"
  exit 1
fi

# start
echo "starting..."
# 获取Jenkins构建号（如果在Jenkins环境中）
BUILD_NUMBER=${BUILD_NUMBER:-"unknown"}
BUILD_HASH=${BUILD_NUMBER}-${GIT_COMMIT}
IMAGE_NAME=select-course:${BUILD_HASH}
PUSH_IMAGE_NAME=swr.cn-north-4.myhuaweicloud.com/bbz/select-course:${BUILD_HASH}



# 构建镜像并使用唯一的标签
docker build -t  "${IMAGE_NAME}" .
echo "build success"

# 推送带有唯一标签的镜像
docker tag "${IMAGE_NAME}"  "${PUSH_IMAGE_NAME}" # 提交带有唯一标签的镜像
docker tag "${IMAGE_NAME}"  select-course:latest # 提交最新的镜像
docker push "${PUSH_IMAGE_NAME}"
docker push select-course:latest
echo "push with unique tag success"

# 清理本地镜像
docker rmi -f "${IMAGE_NAME}"
docker rmi -f select-course:latest
docker rmi -f "${PUSH_IMAGE_NAME}"
echo "clear success"

# k8 replace

# 遍历目录中的所有 YAML 文件
find "./k8s" -name "*.yaml" | while read FILE; do
  sed -i "s|BUILD_HASH|${BUILD_HASH}|" "$FILE"
done
# load configmap
kubectl create configmap config-env -n select-course --from-file=./.env
kubectl create -f k8s
echo "deploy success"
echo "all done"