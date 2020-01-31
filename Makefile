default: clean
	@export GOPATH=$$GOPATH:$$(pwd) && go install agent
run: default
	@bin/agent
	@echo ""
clean:
	@rm -rf bin
	@rm -rf pkg
setup:
	go get gopkg.in/mgo.v2
	go get -u github.com/aws/aws-sdk-go/...
	go get github.com/bsphere/le_go
	go get github.com/bitly/go-simplejson
edit:
	@export GOPATH=$$GOPATH:$$(pwd) && atom .
edit2:
	@export GOPATH=$$GOPATH:$$(pwd) && code .
test:
	@export GOPATH=$$GOPATH:$$(pwd) && go test ./...
test_v:
	@export GOPATH=$$GOPATH:$$(pwd) && go test -v ./...
