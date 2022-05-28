#################################################
GOOS	:= $(shell go env GOOS)
GOARCH	:= $(shell go env GOARCH)
GOFILES	:= $(shell ls *.go |grep -v test)
GOBUILD	:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build


#################################################
default:	deps test build
docker:		deps test linuxbuild docker

deps:
	go get

build:
	CGO_ENABLED=0 $(GOBUILD) $(GOFILES)

linuxbuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFILES)

docker:
	docker build -t maguec/redis-graphql .

test:
	go test -v ./...
