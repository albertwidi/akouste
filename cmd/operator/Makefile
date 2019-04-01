OPERATOR_NAME=appoperator

build:
	go build -v -o $(OPERATOR_NAME) *.go

run:
	make build
	consul-template -config ./template/config/local-template1.hcl \
		- config ./template/config/local-template2.hcl