BIN_PATH=bin
BIN_NAME=sync_wrike_confluence
MAIN_PATH=main.go

.PHONY: run build clean test

run:
	go run ${MAIN_PATH}

build: clean
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -a -ldflags '-s -w' -o ${BIN_PATH}/${BIN_NAME}-linux-amd64 ${MAIN_PATH}
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -a -ldflags '-s -w' -o ${BIN_PATH}/${BIN_NAME}-linux-arm64 ${MAIN_PATH}
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -a -ldflags '-s -w' -o ${BIN_PATH}/${BIN_NAME}-darwin-arm64 ${MAIN_PATH}
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -a -ldflags '-s -w' -o ${BIN_PATH}/${BIN_NAME}-darwin-amd64 ${MAIN_PATH}
	CGO_ENABLED=0 GOOS=windows  GOARCH=amd64 go build -a -ldflags '-s -w' -o ${BIN_PATH}/${BIN_NAME}-windows-amd64.exe ${MAIN_PATH}

clean:
	go clean
	rm -rf ${BIN_PATH}

test:
	${GOPATH}/bin/goconvey
