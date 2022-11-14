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

.PHONY: protoc test build protoc/include

# 任意平台，如果本地有安装 protoc 环境，均可以使用
protoc:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	protoc --proto_path=./api -I ./api/protoc/include/ --go_out=plugins=grpc:./api ./api/*.proto


# Linux 环境下如果没有protoc，可以使用 api/protoc 提供的
PROTOC:=./api/protoc/bin/protoc
protoc/include:
	${PROTOC} --proto_path=./api -I ./api/protoc/include/ --go_out=plugins=grpc:./api ./api/*.proto

# 构建示例
build:
	echo 'client building ....'
	go build -o remote-rate-limit-client examples/client/client.go

	echo 'server-v1 building ...'
	go build -o remote-rate-limit-server-v1 examples/server/v1/server.go

	echo 'server-v2 building ...'
	go build -o remote-rate-limit-server-v2 examples/server/v2/server.go

	echo 'server-v3 building...'
	go build -o remote-rate-limit-server-v3 examples/server/v3/server.go

# 运行示例应用
test: build
	echo 'clean old server...'
	ps aux | grep -v "grep" | grep "rate-limit-server" | awk '{print $$2}' | xargs kill -9

	echo 'server-v3 running in nohup...'
	./remote-rate-limit-server-v3 >> server-v3.log 2>&1 &

	echo 'start client and call plugin server...'
	./remote-rate-limit-client