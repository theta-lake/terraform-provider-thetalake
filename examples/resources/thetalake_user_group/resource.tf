data "thetalake_user" "example" {
  email = "jane.smith@example.com"
}

resource "thetalake_user_group" "example" {
  name        = "Example User Group"
  description = "An example user group managed by Terraform"
  user_ids    = [data.thetalake_user.example.id]
}
