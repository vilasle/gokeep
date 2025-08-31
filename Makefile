BUILD_VERSION=$(shell git describe --always --long)
BUILD_DATE=$(shell date +'%Y/%m/%d %H:%M:%S')
BUILD_COMMIT=$(shell git log --format='%H' -n 1)
server:
	go build -o bin/gokeep-backend -ldflags "-X 'main.buildVersion=$(BUILD_VERSION)' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)'" cmd/server/*.go

client:
	go build -o bin/gokeep-cli -ldflags "-X 'main.buildVersion=$(BUILD_VERSION)' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)'" cmd/client/*.go

clear:
	rm -rf bin/*
	
generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto	

generate-mock:
	mockgen -package=model -destination=internal/model/repository_mock_test.go -source=internal/model/repository.go
	mockgen -package=model -destination=internal/model/encryption_mock_test.go -source=internal/model/encryption.go
	mockgen -package=encryption -destination=internal/encryption/repository_mock_test.go -source=internal/model/repository.go
	
	mockgen -package=auth -destination=internal/service/auth/repository_mock_test.go -source=internal/model/repository.go
	mockgen -package=auth -destination=internal/service/auth/server_repository_mock_test.go -source=internal/repository/server/repository.go
	mockgen -package=auth -destination=internal/service/auth/user_mock_test.go -source=internal/model/user.go
	
	mockgen -package=private -destination=internal/service/private/repository_mock_test.go -source=internal/model/repository.go
	mockgen -package=private -destination=internal/service/private/encryption_mock_test.go -source=internal/encryption/encryption.go
	mockgen -package=private -destination=internal/service/private/user_mock_test.go -source=internal/model/user.go
	
	mockgen -package=server -destination=internal/server/service_mock_test.go -source=internal/service/service.go

	mockgen -package=grpc -destination=internal/service/client/grpc/pb_mock_test.go -source=proto/service_grpc.pb.go

	mockgen -package=cli -destination=internal/client/cli/service_mock_test.go -source=internal/service/client/service.go
	mockgen -package=cli -destination=internal/client/cli/repository_mock_test.go -source=internal/repository/client/repository.go
	mockgen -package=cli -destination=internal/client/cli/encryption_mock_test.go -source=internal/encryption/encryption.go
	

test:
	go test ./...

coverage:
	go test ./... -coverprofile cover.out
	go tool cover -html cover.out
	rm cover.out

coverage-percent:
	go test ./... -coverprofile=cover.out
	go tool cover -func cover.out | tail -n 1 && rm -rf cover.out
	