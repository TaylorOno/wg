.PHONY: sonar-scanner
## sonar: *REQUIRES LOCAL SONARQUBE and SONAR-SCANNER run uploads to local sonar instance
sonar: sonar-scanner
	sonar-scanner -Dsonar.projectKey=trello-tribbles -Dsonar.exclusions=**/*_test.go,**/test_data/*,**/mocks/**,**/main.go -Dsonar.host.url=http://localhost:9000 -Dsonar.source=. -Dsonar.go.coverage.reportPaths=**/coverage.out

.PHONY: test
## test: runs all tests
test:
	go test ./... -coverprofile=./coverage.out

.PHONY: vet
## vet: runs go vet
vet:
	go vet ./...

.PHONY: fmt
## fmt: runs go fmt
fmt:
	go fmt ./...

.PHONY: pre-release
## pre-release: runs all tests and go tools
pre-release: fmt vet test

.PHONE: build
## Builds the binary for release
build:
	go build -ldflags='-w -s -extldflags' -a -o azure_function/bin/wg ./cmd

.PHONE: run
## Runs the application
run:
	go run ./cmd

.PHONE: start
## Start local development instance using azure function runtime
start: build
	cd azure_function && func start

.PHONE: deploy
## Deploys the function to an existing functionapp
deploy: pre-release
	set CGO_ENABLED=0
	set GOOS=windows
	set GOARCH=amd64
	go build -ldflags='-w -s -extldflags' -a -o azure_function/bin/wg ./cmd
	upx --brute -qq azure_function/bin/wg
	cd azure_function && func azure functionapp publish az-wgdrinks

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/ /'
