OPERATOR_NAME=appoperator
OPERATOR_DIR=./cmd/operator
DOWNLOADER_NAME=configdownloader
DOWNLOADER_DIR=./cmd/downloader

build-operator:
	go build -v -o $(OPERATOR_NAME) $(OPERATOR_DIR)/*.go

build-downloader:
	go build -v -o $(DOWNLOADER_NAME) ./cmd/downloader/*.go

docker-build-operator:
	cd $(OPERATOR_DIR) && go build -v -o $(OPERATOR_NAME) *.go
	cd $(OPERATOR_DIR) && docker build .

run-opeator:
	make build-operator
	./$(OPERATOR_NAME)

run-downloader:
	make build-downloader
	./$(DOWNLOADER_NAME)