#!/bin/bash

IMAGEHUB="registry.cn-shenzhen.aliyuncs.com/huanglin_hub"
PROJECT="openim-msggateway-proxy"
ALLPRO="all"

echo "building ${PROJECT}"
DOCKER_PUSHIMG=${IMAGEHUB}/${PROJECT}:dev
docker rmi  ${DOCKER_PUSHIMG}
docker build -f ./Dockerfile -t ${DOCKER_PUSHIMG} .
docker push ${DOCKER_PUSHIMG}