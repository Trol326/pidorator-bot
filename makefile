start:
	go run main.go

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: lint
lint:
	golangci-lint run --enable-all

.PHONY: install-goctl	
install-goctl:
	go install github.com/zeromicro/go-zero/tools/goctl@latest

.PHONY: gen-dockerfile
gen-dockerfile:
	goctl docker -go app/main.go -tz "Europe/Moscow"

.PHONY: install-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
	golangci-lint --version

.PHONY: install
install:
	tidy install-goctl install-lint