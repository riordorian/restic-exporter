include .env

.SILENT:

convert_proto:
	protoc --proto_path=./pkg/proto \
 		--go_out=./pkg/proto_gen/ \
 		--go-grpc_out=./pkg/proto_gen/ \
 		--grpc-gateway_out=./pkg/proto_gen/ \
 		--openapiv2_out ./internal/infrastructure/ports/http/api/ \
 		--openapiv2_opt ignore_comments=true \
 		--openapiv2_opt allow_merge=true \
 		--openapiv2_opt generate_unbound_methods=false \
 		./pkg/proto/*.proto && \
 	cd pkg/proto_gen/grpc && protoc-go-inject-tag -input="*.pb.go"

build_frontend:
		mkdir -p ./frontend/src/proto
		protoc $(find ./pkg/proto -iname "*.proto") \
 			--proto_path=./pkg/proto \
			--plugin=protoc-gen-grpc-web=./frontend/node_modules/.bin/protoc-gen-grpc-web \
			--js_out=import_style=commonjs:./frontend/src/proto \
			--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./frontend/src/proto

