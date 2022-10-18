BINARY_NAME=goprivate
 
all: clean build test
 
build: webfront
	go build -o ${BINARY_NAME} cmd/goprivate/main.go
 
test:
	go test ./...
 
run: build
	./${BINARY_NAME}

.PHONY: webfront
webfront:
	./compile-webfront.sh

clean:
	rm -rf webfront/dist/
	rm -rf server/webfront
	rm -f ${BINARY_NAME}