all: build install 

build:
	go get gopkg.in/yaml.v2 && go build

install:
	go install

test:
	go test -v *.go
