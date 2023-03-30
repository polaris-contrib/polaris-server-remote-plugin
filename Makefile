# Tencent is pleased to support the open source community by making Polaris available.
#
# Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
#
# Licensed under the BSD 3-Clause License (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# https://opensource.org/licenses/BSD-3-Clause
#
# Unless required by applicable law or agreed to in writing, software distributed
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.

.PHONY: protoc test build protoc/*

# 任意平台，如果本地有安装 protoc 环境，均可以使用
protoc:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --proto_path=./api -I ./api/protoc/include/ --go-grpc_out=require_unimplemented_servers=false:./api --plugin=protoc-gen-grpc  ./api/*.proto
	protoc --proto_path=./api -I ./api/protoc/include/ --go_out=./api ./api/*.proto


# Linux 环境下如果没有protoc，可以使用 api/protoc 提供的
PROTOC:=./api/protoc/bin/protoc
protoc/include:
	${PROTOC} --proto_path=./api -I ./api/protoc/include/ --go-grpc_out=require_unimplemented_servers=false:./api --plugin=protoc-gen-grpc  ./api/*.proto
	${PROTOC} --proto_path=./api -I ./api/protoc/include/ --go_out=./api ./api/*.proto

# python grpc
protoc/python:
	pip3 install grpcio-tools
	python3 -m grpc_tools.protoc \
		--proto_path=./api -I ./api/protoc/include/ \
		--python_out=./examples/server/v4 \
		--pyi_out=./examples/server/v4 \
		--grpc_python_out=./examples/server/v4 \
		./api/*.proto

# 构建 Python 环境
build/python:
	pip3 install pyinstaller
	pip3 install grpcio-tools
	pyinstaller --specpath=./examples/server/v4 --distpath=./ -n=remote-rate-limit-server-v4 -F ./examples/server/v4/server.py