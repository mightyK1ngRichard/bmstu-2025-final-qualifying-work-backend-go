run:
	go run cmd/cake/main.go

tests:
	go test -v -p 2 \
	./internal/pkg/auth/delivery/grpc \
	./internal/pkg/auth/usecase

db_restart:
	docker compose up -d

auth_proto:
	cd proto && \
	protoc --go_out=../internal/pkg/auth/delivery/grpc/generated \
	       --go-grpc_out=../internal/pkg/auth/delivery/grpc/generated \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       auth.proto && \
	protoc --swift_out=../internal/pkg/auth/delivery/grpc/generated \
	       --grpc-swift_out=../internal/pkg/auth/delivery/grpc/generated \
	       auth.proto

cake_proto:
	cd proto && \
	protoc --go_out=../internal/pkg/cake/delivery/grpc/generated \
	       --go-grpc_out=../internal/pkg/cake/delivery/grpc/generated \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       cake.proto && \
	protoc --swift_out=../internal/pkg/cake/delivery/grpc/generated \
	       --grpc-swift_out=../internal/pkg/cake/delivery/grpc/generated \
	       cake.proto

profile_proto:
	cd proto && \
	protoc --go_out=../internal/pkg/profile/delivery/grpc/generated \
	       --go-grpc_out=../internal/pkg/profile/delivery/grpc/generated \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       profile.proto && \
	protoc --swift_out=../internal/pkg/profile/delivery/grpc/generated \
	       --grpc-swift_out=../internal/pkg/profile/delivery/grpc/generated \
	       profile.proto