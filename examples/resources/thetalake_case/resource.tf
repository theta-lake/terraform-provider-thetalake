data "thetalake_user" "example" {
  email = "user@example.com"
}

resource "thetalake_case" "example" {
  name        = "Example Case"
  number      = "CASE-001"
  description = "An example case managed by Terraform"
  open_date   = "2024-01-15T10:00:00Z"
  visibility  = "PRIVATE"

  manager_ids = [data.thetalake_user.example.id]
}
