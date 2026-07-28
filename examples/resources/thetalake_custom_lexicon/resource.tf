resource "thetalake_custom_lexicon" "example" {
  name        = "Example Custom Lexicon"
  description = "Detects mentions of confidential project code names"
  risk_type   = "risk"

  rules = [
    "project-nightingale",
    "project-falcon",
  ]

  rule_scope              = ["chat", "email", "doc"]
  communication_direction = ["inbound", "outbound"]
  attachments_enabled     = true
  filename_analyzed       = true

  policy_ids = [1, 2]

  start_date = "2024-01-01"
}
