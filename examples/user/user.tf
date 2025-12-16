terraform {
    required_providers {
        thetalake = {
            source = "registry.terraform.io/thetalake/thetalake"
            version = "0.1.0"
        }
    }
}

provider "thetalake" {
    endpoint = "insert-endpoint"
    token = "insert-api-token"
}

resource "thetalake_user" "user_01" {
    name = "insert-name"
    email = "insert-email@thetalake.com"
    password = "insert-password"
    password_confirmation = "insert-password"
    role_id = 3
}

output "user_id" {
    value = thetalake_user.user_01.id
}