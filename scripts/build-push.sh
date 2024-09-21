#!/bin/bash
# clear old image
docker rmi -f select-course:latest
# build
docker build -t select-course:latest .
# push
sudo docker tag select-course:latest swr.cn-north-4.myhuaweicloud.com/bbz/select-course:latest
sudo docker push swr.cn-north-4.myhuaweicloud.com/bbz/select-course:latest
echo "push success"

# delete
docker rmi -f select-course:latest
