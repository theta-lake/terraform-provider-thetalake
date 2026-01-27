build:
	go build -o terraform-provider-thetalake .

install:
	@go install .

docs:
	@tfplugindocs generate