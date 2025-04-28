run:
	go run cmd/cake/main.go & \
    	go run cmd/profile/main.go & \
    	go run cmd/auth/main.go & \
    	go run cmd/chat/main.go & \
    	go run cmd/reviews/main.go & \
    	go run cmd/order/main.go & \
    	wait

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

chat_proto:
	cd proto && \
	protoc --go_out=../internal/pkg/chat/delivery/grpc/generated \
	       --go-grpc_out=../internal/pkg/chat/delivery/grpc/generated \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       chat.proto && \
	protoc --swift_out=../internal/pkg/chat/delivery/grpc/generated \
	       --grpc-swift_out=../internal/pkg/chat/delivery/grpc/generated \
	       chat.proto

feedback_proto:
	cd proto && \
	protoc --go_out=../internal/pkg/reviews/delivery/grpc/generated \
	       --go-grpc_out=../internal/pkg/reviews/delivery/grpc/generated \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       feedback.proto && \
	protoc --swift_out=../internal/pkg/reviews/delivery/grpc/generated \
	       --grpc-swift_out=../internal/pkg/reviews/delivery/grpc/generated \
	       feedback.proto

order_proto:
	cd proto && \
		protoc --go_out=../internal/pkg/order/delivery/grpc/generated \
			   --go-grpc_out=../internal/pkg/order/delivery/grpc/generated \
			   --go_opt=paths=source_relative \
			   --go-grpc_opt=paths=source_relative \
			   order.proto && \
		protoc --swift_out=../internal/pkg/order/delivery/grpc/generated \
			   --grpc-swift_out=../internal/pkg/order/delivery/grpc/generated \
			   order.proto
