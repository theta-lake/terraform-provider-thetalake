data "thetalake_user" "alice" {
  email = "alice@example.com"
}

data "thetalake_user" "bob" {
  email = "bob@example.com"
}

# Workspaces cannot be created via the API, so an existing workspace must be
# imported before Terraform can manage it. You can do this with:
# `terraform import thetalake_workspace.example id`, where id is the ID of the
# workspace to import.
resource "thetalake_workspace" "example" {
  allow_anonymous_via_shared_links   = false
  analysis_supervision_space_ids     = [1]
  audit_log_retention_period         = 365
  case_management_manager_assignment = false
  default_transcription_language     = "en"
  default_workspace_timezone         = "Etc/UTC"
  delete_on_expiration               = false
  hide_attachments_from_search       = false
  preferred_languages                = ["en", "es"]
  reauthenticate_on_network_change   = false
  shared_links_expiration_period     = 7
  show_system_messages_in_chat       = false
  use_name_matcher                   = true
  use_owner_only_space_matcher       = false
  user_ids = [
    data.thetalake_user.alice.id,
    data.thetalake_user.bob.id,
  ]
}