#!/bin/bash

# 构建镜像
docker build -t select-course .
# 启动服务
docker-compose up --build
