base_dir = `pwd`
gopath = "$(base_dir)/third_party:$(GOPATH)"

all: clean fmt test
	@echo
	@echo "==> Compiling source code."
	@echo
	@env GOPATH=$(gopath) go build -v -o ./bin/sift ./sift
	@env GOPATH=$(gopath) go build -v -o ./bin/sift-source-aws ./source/aws
	@env GOPATH=$(gopath) go build -v -o ./bin/sift-source-test-aws ./source/test/aws
	@chmod a+x ./bin/*
test:
	@echo
	@echo "==> Running tests."
	@echo
	@env GOPATH=$(gopath) go test $(COVER) ./sift/...
deps:
	@echo
	@echo "==> Downloading dependencies."
	@echo
	@env GOPATH=$(gopath) go get -d -v ./...
	@echo
	@echo "==> Removing .git, .bzr, and .hg from third_party."
	@echo
	@find ./third_party -type d -name .git | xargs rm -rf
	@find ./third_party -type d -name .bzr | xargs rm -rf
	@find ./third_party -type d -name .hg | xargs rm -rf
fmt:
	@echo
	@echo "==> Formatting source code."
	@echo
	@gofmt -w -tabs=false -tabwidth=2 ./sift ./source
clean:
	@echo
	@echo "==> Cleaning up previous builds."
	@echo
	@rm -rf bin pkg third_party/pkg

.PHONY: all clean deps format test
