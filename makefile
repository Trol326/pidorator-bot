start:
	go run app/main.go

.PHONY: tidy
tidy:
	go get github.com/bwmarrin/discordgo@master
	go mod tidy -v

.PHONY: lint
lint:
	golangci-lint run 

.PHONY: install-goctl	
install-goctl:
	go install github.com/zeromicro/go-zero/tools/goctl@latest

.PHONY: gen-dockerfile
gen-dockerfile:
	goctl docker -go app/main.go -tz "Etc/GMT"

.PHONY: build-image
build-image:
	docker build -t bot:v1 .

db/up:
	docker-compose -f docker-compose-db.dev.yml up -d

db/down:
	docker-compose -f docker-compose-db.dev.yml down

.PHONY: install-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
	golangci-lint --version

.PHONY: install
install: tidy install-goctl install-lint

check: tidy lint