#!/bin/bash

if [ ! "$(basename "$PWD")" ] = "select-course"; then
  echo "please run this script in select-course"
  exit 1
fi

# start
echo "starting..."
# 获取Jenkins构建号（如果在Jenkins环境中）
BUILD_NUMBER=${BUILD_NUMBER:-"unknown"}
IMAGE_NAME=select-course:${BUILD_NUMBER}-${GIT_COMMIT}
PUSH_IMAGE_NAME=swr.cn-north-4.myhuaweicloud.com/bbz/select-course:${BUILD_NUMBER}-${GIT_COMMIT}


# 构建镜像并使用唯一的标签
docker build -t  "${IMAGE_NAME}" .
echo "build success"

# 推送带有唯一标签的镜像
docker tag "${IMAGE_NAME}"  "${PUSH_IMAGE_NAME}"
docker push "${PUSH_IMAGE_NAME}"
echo "push with unique tag success"

# 清理本地镜像
docker rmi -f "${IMAGE_NAME}"
echo "clear success"

# k8 replace

# load configmap
kubectl replace configmap config-env -n select-course --from-file=./.env
kubectl replace -f k8s
echo "replace success"