resource "thetalake_role" "example" {
  name        = "Example Role"
  description = "Permissions for a user tasked solely with reviewing content"
  permissions = [
    "cases:create",
    "cases:read",
    "cases:update",
    "cases:delete",
  ]
}
