# terraform-provider-thetalake
Terraform Provider for Theta Lake

[![CI](https://github.com/theta-lake/terraform-provider-thetalake/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/theta-lake/terraform-provider-thetalake/actions/workflows/ci.yml)

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

4. Create and cd into a test directory
```bash
mkdir provider-test

# cd into the directory
cd provider-test
```

5. Create your `main.tf` file and save it in the test directory
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

6. Run Terraform plan
```bash
# terraform init is not used for local dev testing
terraform plan
```

7. Run Terraform apply
```bash
terraform apply
```

## Generate Documentation
1. Install the tool
```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# verify the installation
tfplugindocs --version
```

2. Generate the docs
```bash
# make sure to run from the root of the repo
make generate-docs
```
