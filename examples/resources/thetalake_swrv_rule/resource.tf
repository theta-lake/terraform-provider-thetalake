resource "thetalake_swrv_rule" "example" {
  name                 = "customer-escalations"
  description          = "Review customer escalations routed into the workflow"
  policy_id            = 147
  retention_library_id = 1
  workflow_id          = 14536
  priority             = 4

  input_sources = [{
    type = "all_uploads"
  }]
}