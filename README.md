# Theta Lake Provider for Terraform

The Theta Lake Provider enables you to manage Theta Lake compliance, archiving, supervision, and risk detection configurations using Terraform’s Infrastructure as Code workflows.

## Requirements

The Theta Lake Provider requires [Terraform](https://developer.hashicorp.com/terraform) 1.13.x or newer.

The provider has been tested with Terraform up to version 1.13.4. Versions newer than 1.13.4 may work, but are not officially supported.

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
    required_providers {
        thetalake = {
            source  = "thetalake/thetalake"
            version = "~> 1.0"
        }
    }
}

provider "thetalake" {
    api_server    = "API URL here"
    client_id     = "client-id"
    client_secret = "client-secret"
}
```

Then run

```bash
terraform init
```

## Documentation

You can find documentation for the Theta Lake Provider on the [Terraform Registry](https://registry.terraform.io/providers/thetalake/thetalake/latest) site.

## Authentication Credentials

You can find documentation on obtaining Theta Lake OAuth credentials at the [Theta Lake Developer](https://developer.thetalake.ai/documentation/api-keys) site.

## License

Copyright 2026 Theta Lake, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
