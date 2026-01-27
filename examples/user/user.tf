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

variable "default_user_password" {
  type      = string
  default   = "Testtesttest123"
  sensitive = true
}

data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

data "thetalake_role" "api_only" {
  name = "API Only"
}

resource "thetalake_user" "user_01" {
  name     = "insert-name"
  email    = "insert-email@thetalake.com"
  password = var.default_user_password
  disabled = false
  role_id  = data.thetalake_role.reviewer.id
}