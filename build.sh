#!/bin/bash

PS3="请选择要编译的系统环境："

options=("Windows_amd64" "linux_amd64")

select action in "${options[@]}"
do
  case $action in
    "Windows_amd64")
      echo "编译Windows版本64位"
      export CGO_ENABLED=0
      export GOOS=windows
      export GOARCH=amd64
      go build -o project-user/target/project-user.exe project-user/main.go
      go build -o project-api/target/project-api.exe project-api/main.go
      break
      ;;
    "linux_amd64")
      echo "编译Linux版本64位"
      export CGO_ENABLED=0
      export GOOS=linux
      export GOARCH=amd64
      go build -o project-user/target/project-user project-user/main.go
      go build -o project-api/target/project-api project-api/main.go
      break
      ;;
    *)
      echo "无效的选择，请重新选择."
      ;;
  esac
done