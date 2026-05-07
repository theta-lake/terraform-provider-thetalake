.PHONY: build install docs test

build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

docs:
	terraform fmt -recursive ./examples/
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.24.0 generate --provider-dir . -provider-name thetalake

test:
	go test $(shell go list ./internal/... | grep -v acctest) -json -coverprofile=coverage.out | go run github.com/mfridman/tparse@latest -notests -follow