#!/bin/bash
# 是否在select-course目录下
if [ ! -d "select-course" ]; then
    echo "请进入select-course目录后再执行脚本"
    exit 1
fi
# 构建镜像
docker build -t select-course .
# 启动服务
docker-compose up --build
