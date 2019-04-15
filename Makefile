.PHONY: clean all

OPERATOR_NAME   := appoperator
OPERATOR_DIR    := ./cmd/operator
DOWNLOADER_NAME := configdownloader
DOWNLOADER_DIR  := ./cmd/downloader
DOWNLOADER_BIN  := ${DOWNLOADER_DIR}/build

.PHONY: test
test:
	go test -race -v ./...

.PHONY: build-operator
build-operator:
	go build -v -o $(OPERATOR_NAME) $(OPERATOR_DIR)/*.go

.PHONY: run-operator
run-opeator:
	make build-operator
	./$(OPERATOR_NAME)

.PHONY: docker-build-operator
docker-build-operator:
	cd $(OPERATOR_DIR) && go build -v -o $(OPERATOR_NAME) *.go
	cd $(OPERATOR_DIR) && docker build .

.PHONY: build-downloader
build-downloader:
	@GO111MODULE=off go build -v -o $(DOWNLOADER_BIN)/$(DOWNLOADER_NAME) $(DOWNLOADER_DIR)/*.go

.PHONY: run-downloader
run-downloader:
	@make build-downloader
	./$(DOWNLOADER_BIN)/$(DOWNLOADER_NAME) \
		-logLevel debug \
		-bucketProto "local" \
		-bucketName "test/local-bucket" \
		-downloadDIR "test/local-downloads" \
		-keepOldCount 2

.PHONY: docker-build-downloader
docker-build-downloader:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build-downloader
	@cd $(DOWNLOADER_BIN) && docker build -t $(DOWNLOADER_NAME) .

.PHONY: docker-run-downloader
docker-run-downloader:
	docker run --env DOCKER_TEST=true --name $(DOWNLOADER_NAME) -dit $(DOWNLOADER_NAME)
	@docker exec -it $(DOWNLOADER_NAME) /bin/bash
	@make docker-stop-downloader

.PHONY: docker-stop-downloader
docker-stop-downloader:
	@docker rm -f $(DOWNLOADER_NAME) 2>/dev/null || true
