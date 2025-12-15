# terraform-provider-thetalake
Terraform Provider for Theta Lake

## Local Testing
1. Install the provider
```bash
# Installs the provider to the bin directory of your go path
make install
```

2. Create a `~/.terraformrc` file if it does not already exist

3. Update your `~/.terraformrc` file
```txt
# the direct section is needed to ensure that the other providers run properly
provider_installation {
    dev_overrides {
        "registry.terraform.io/thetalake/thetalake" = "/Users/tyang/go/bin"
    }
    direct {
        enabled = true
    }
}
```

4. Create your `main.tf` file
```txt
terraform {
    required_providers {
        thetalake = {
            source = "registry.terraform.io/thetalake/thetalake"
            version = "0.1.0"
        }
    }
}

provider "thetalake" {
    endpoint = "http://localhost:6002/api/v1/users"
    token = "insert-token"
}

resource "thetalake_user" "user_01" {
    name = "insert-name"
    email = "insert-email@thetalake.com"
    password = "insert-password"
    password_confirmation = "insert-password"
    role_id = 3
    search_id = 372899
}
```

5. Run `terraform plan`. You cannot run `terraform init` since this is a local dev test