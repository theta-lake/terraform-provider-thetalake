.PHONY: build install docs

build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

docs:
	tfplugindocs generate --provider-name thetalake