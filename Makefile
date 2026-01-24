build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

generate-docs:
	@tfplugindocs generate