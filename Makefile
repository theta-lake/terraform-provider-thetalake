.PHONY: build install docs test testacc testall

-include local.env
export

build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

vulncheck:
	govulncheck ./... 

docs:
	terraform fmt -recursive ./examples/
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.24.0 generate --provider-dir . -provider-name thetalake

test:
	go test $(shell go list ./internal/... | grep -v acctest) -json -coverprofile=coverage.out | go run github.com/mfridman/tparse@latest -notests -follow

testacc:
	TF_ACC=1 go test ./internal/... -run TestAcc -v -count=1

testall:
	TF_ACC=1 go test $(shell go list ./internal/... | grep -v acctest) -count=1 -p 1 -json -coverprofile=coverage.out | go run github.com/mfridman/tparse@latest -notests -follow