#!/bin/bash

if [ ! "$(basename "$PWD")" ] = "select-course"; then
  echo "please run this script in select-course"
  exit 1
fi

# start
echo "starting..."

# 获取当前Git提交的前7位
GIT_COMMIT_SHORT=$(git rev-parse --short HEAD)

# 获取Jenkins构建号（如果在Jenkins环境中）
BUILD_NUMBER=${BUILD_NUMBER:-"unknown"}

# 构建镜像并使用唯一的标签
docker build -t select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT} .

echo "build success"

# 推送带有唯一标签的镜像
docker tag select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT} swr.cn-north-4.myhuaweicloud.com/bbz/select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT}
docker push swr.cn-north-4.myhuaweicloud.com/bbz/select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT}

echo "push with unique tag success"

# 如果是最后一个版本，额外添加 latest 标签并推送
if [ "$BUILD_NUMBER" = "last" ]; then
  # 重新打上 latest 标签
  docker tag select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT} swr.cn-north-4.myhuaweicloud.com/bbz/select-course:latest
  # 推送 latest 标签
  docker push swr.cn-north-4.myhuaweicloud.com/bbz/select-course:latest
  echo "push with latest tag success"
fi

# 清理本地镜像
docker rmi -f select-course:${BUILD_NUMBER}-${GIT_COMMIT_SHORT}

echo "clear success"

# k8 replace
kubectl replace -f k8s
echo "replace success"