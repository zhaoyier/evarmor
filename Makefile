proto:
	protoc -I apidoc apidoc/proto/common/*.proto --go_out=plugins=grpc:rpc/
	protoc -I apidoc apidoc/proto/chat/*.proto --go_out=plugins=grpc:rpc/

clean:
	find . -name "*.pb.go"|xargs rm -rf
