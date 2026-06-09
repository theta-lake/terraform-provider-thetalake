resource "thetalake_retention_library" "example" {
  name               = "Example Retention Library"
  storage_account_id = 1
  description        = "Used to store data for the compliance environment"
  external_id        = "ext-001"

  retention_period_enabled = true
  retention_period_days    = 365
  retain_in_review         = true

  sec_compliant_storage_enabled = false
}
