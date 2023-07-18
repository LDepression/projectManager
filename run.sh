#!/bin/bash

# 设置字符编码
chcp 65001

# 进入项目目录
cd project-user || exit

# 构建 Docker 镜像
docker build -t project-user:latest .

# 返回上级目录
cd ..

# 启动 Docker 容器
docker-compose up -d