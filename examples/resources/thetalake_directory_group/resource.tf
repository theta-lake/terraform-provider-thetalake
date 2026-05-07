data "thetalake_identity" "example" {
  email = "user@example.com"
}

resource "thetalake_directory_group" "example" {
  name         = "Example Directory Group"
  description  = "An example directory group managed by Terraform"
  external_id  = "dg-ext-001"
  identity_ids = [data.thetalake_identity.example.id]
}
