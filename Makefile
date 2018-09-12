GOPATH := ${PWD}/vendor:${GOPATH}
build:
	go build -o ${BINARY_NAME} -a -v -ldflags '-extldflags "--static"' ./src/main

clean:
	rm -rf ${BINARY_NAME}