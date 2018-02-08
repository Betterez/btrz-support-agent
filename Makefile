default:
	@export GOPATH=$$GOPATH:$$(pwd) && go install agent
run: default
	@bin/agent
	@echo ""
clean:
	@rm -rf bin
setup:
	go get gopkg.in/mgo.v2
	go get -u github.com/aws/aws-sdk-go/...
	go get github.com/bsphere/le_go
edit:
	@export GOPATH=$$GOPATH:$$(pwd) && atom .
test:
	@export GOPATH=$$GOPATH:$$(pwd) && go test ./...
test_v:
	@export GOPATH=$$GOPATH:$$(pwd) && go test -v ./...
