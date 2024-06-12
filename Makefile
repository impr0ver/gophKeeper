CLIENTNAME=client
SERVERNAME=server
BUILDFLAGS=-ldflags="-X 'main.buildVersion=v1.00' -X 'main.buildDate=$(shell date -u +'%Y-%m-%d %H:%M:%S')' -X 'main.buildCommit=${shell git rev-parse HEAD}'"

cert:
	cd cmd/cert; ./gen.sh; cd ../..;

build_client:
	GOARCH=amd64 GOOS=darwin go build -o ./cmd/client/${CLIENTNAME}-darwin ${BUILDFLAGS} cmd/client/main.go
	GOARCH=amd64 GOOS=windows go build -o ./cmd/client/${CLIENTNAME}-windows ${BUILDFLAGS} cmd/client/main.go
	GOARCH=amd64 GOOS=linux go build -o ./cmd/client/${CLIENTNAME}-linux ${BUILDFLAGS} cmd/client/main.go
	
build_server:
	GOARCH=amd64 GOOS=darwin go build -o ./cmd/server/${SERVERNAME}-darwin ${BUILDFLAGS} cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ./cmd/server/${SERVERNAME}-windows ${BUILDFLAGS} cmd/server/main.go
	GOARCH=amd64 GOOS=linux go build -o ./cmd/server/${SERVERNAME}-linux ${BUILDFLAGS} cmd/server/main.go
	
build:
	go build -o ./cmd/client/${CLIENTNAME} ${BUILDFLAGS} cmd/client/main.go
	go build -o ./cmd/server/${SERVERNAME} ${BUILDFLAGS} cmd/server/main.go
	
run_client:
	go run ./cmd/client/${CLIENTNAME}

run_server:
	go run ./cmd/server/${SERVERNAME}

test:
	go clean -testcache
	go test ./...

cov:
	go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out | grep total | awk '{print $3}'

clean:
	go clean
	rm ./cmd/server/${SERVERNAME}-darwin ./cmd/server/${SERVERNAME}-linux ./cmd/server/${SERVERNAME}-windows
	rm ./cmd/client/${CLIENTNAME}-darwin ./cmd/client/${CLIENTNAME}-linux ./bcmd/clientin/${CLIENTNAME}-windows
	rm ./cmd/server/${SERVERNAME} ./cmd/client/${CLIENTNAME}

