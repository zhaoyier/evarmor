export PATH := $(shell pwd)/node_modules/.bin:$(PATH)
SHELL := /bin/bash


# 项目初始化
init:
	yarn
	git submodule init
	git submodule update

# 开发模式
dev:
	npm start

build:clean
	npm run build

publish:clean
	yarn
	npm run publish

clean:
	rm -rf dist

# 建议每次只生成自己的服务
# 多个 service 文件以空格隔开
thriftServices = apidoc/thrift/EzSeller.thrift
protobufServices = 

genservices:
	@$(foreach var, $(protobufServices), protoc --plugin=protoc-gen-json-ts=./node_modules/protoc-gen-json-ts/bin/protoc-gen-json-ts --json-ts_out=:src/services -I ./apidoc/proto $(var);)
	@$(foreach var, $(thriftServices), tgen gen -l typescript -m rest -i $(var) -o ./src/services;)
