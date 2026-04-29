data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

resource "thetalake_user" "example" {
  name     = "Jane Smith"
  email    = "jane.smith@example.com"
  password = "ExamplePassword123"
  disabled = false
  role_id  = data.thetalake_role.reviewer.id
}
