data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

data "thetalake_integration" "example" {
  name = "Example Integration"
}

data "thetalake_retention_library" "example" {
  name = "Example Retention Library"
}

data "thetalake_directory_group" "example" {
  name = "Example Directory Group"
}

data "thetalake_user_group" "example" {
  name = "Example User Group"
}

resource "thetalake_user" "example" {
  name     = "Jane Smith"
  email    = "jane.smith@example.com"
  password = "ExamplePassword123"
  disabled = false
  role_id  = data.thetalake_role.reviewer.id
}

resource "thetalake_supervision_space" "example" {
  all_participants                     = false
  all_users                            = false
  name                                 = "Example Supervision Space"
  description                          = "An example supervision space managed by Terraform"
  directory_group_ids                  = [data.thetalake_directory_group.example.id]
  external_id                          = "space-ext-001"
  hard_enforce                         = false
  integration_ids                      = [data.thetalake_integration.example.id]
  media_types                          = ["chat", "email"]
  retention_library_ids                = [data.thetalake_retention_library.example.id]
  requested_supervision_space_priority = 100
  user_group_ids                       = [data.thetalake_user_group.example.id]
  user_ids                             = [resource.thetalake_user.example.id]
}
