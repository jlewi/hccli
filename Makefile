tidy:
	gofmt -s -w .
	goimports -w .

lint:
	# golangci-lint automatically searches up the root tree for configuration files.
	golangci-lint run

build:
	CGO_ENABLED=0 go build  github.com/jlewi/hccli
