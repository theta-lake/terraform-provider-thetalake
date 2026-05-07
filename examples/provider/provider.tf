terraform {
  required_version = ">= 1.13.0"
  required_providers {
    thetalake = {
      source  = "registry.terraform.io/thetalake/thetalake"
      version = "0.1.0"
    }
  }
}

provider "thetalake" {
  api_server    = "insert-api-server"
  client_id     = "insert-api-token"
  client_secret = "insert-api-secret"
}
