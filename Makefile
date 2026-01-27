.PHONY: build install docs

build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

docs:
	go generate ./tools