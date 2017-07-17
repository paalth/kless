protoc -I pkg/klessserver/grpc/ pkg/klessserver/grpc/klessserver.proto --go_out=plugins=grpc:pkg/klessserver/grpc
