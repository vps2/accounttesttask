.PHONY: protoc test run-server run-client build-server build-client build clean help
#.SILENT:

GOOS = $(shell go env GOOS)
ADDR?=:8080
POLLING_INTERVAL?=15s
IDLE?=15s

## protoс: сгенерировать go файл по описанию
protoc:
	protoc -I api/proto --go_out=plugins=grpc:internal/api api/proto/accounts.proto
	protoc -I api/proto --go_out=plugins=grpc:internal/api api/proto/statistics.proto

## build-server: создать исполняемый файл сервера
build-server:
ifeq ($(GOOS),windows)
	go mod download && go build -o bin/server.exe -ldflags "-s -w" cmd/server/main.go
else
	#go build -o bin/server -ldflags "-s -w" cmd/server/main.go
	#компилируем для запуска в docker
	go mod download && CGO_ENABLED=0 go build -o bin/server -ldflags "-s -w" cmd/server/main.go
	#делаем образ docker
	docker build -t accounts-srv -f Dockerfile .
endif

## build-client: создать исполняемый файл клиента
build-client:
ifeq ($(GOOS),windows)
	go mod download && go build -o bin/client.exe -ldflags "-s -w" cmd/client/main.go
else
	go mod download && go build -o bin/client -ldflags "-s -w" cmd/client/main.go
endif

## build : создать исполняемые файл сервера и клиента
build: build-server build-client

## run-server: запустить сервер. Можно установить значения переменной ADDR
run-server: build-server
ifeq ($(GOOS),windows)
	go run cmd/server/main.go --addr $(ADDR) --polling-inteval $(POLLING_INTERVAL)
else
	docker run --rm --name accounts-srv -p $(ADDR):8080 -e POLLING_INTERVAL=$(POLLING_INTERVAL) accounts-srv
endif

## run-client: запустить клиента. Можно установить значения переменной IDLE
run-client: build-client
	go run cmd/client/main.go -idle $(IDLE) -cfg-file "configs/config.yml"

## test  : запустить тестирование
test:
	go test ./...

## clean : удалить содержимое папки bin
clean:
ifeq ($(GOOS),windows)
	powershell "Get-ChildItem bin/* -Recurse | Remove-Item -Recurse"
else
	rm -rf bin/*
endif

help: Makefile
ifeq ($(GOOS),windows)
	@powershell '(Get-Content $< -Encoding utf8) -match "^##" -replace "^##(.*?):\s(.*?)"," `$$1`t`$$2"'
else
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
endif