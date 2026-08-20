variable "undeliverable_mailbox_password" {
  type        = string
  sensitive   = true
  description = "Password for the mailbox that receives undeliverable email notifications."
}

resource "thetalake_integration" "journaling" {
  name   = "Custom Generic Journaling Integration"
  paused = false

  generic_journaling = {
    download_o365_onedrive_links = true
    index_headers                = "X-Header-Score,X-Routed-Via"
    undeliverable_email_server   = "email.example.com"
    undeliverable_email_user     = "Undeliverable User"
    undeliverable_email_password = var.undeliverable_mailbox_password
    undeliverable_email_port     = 993
  }
}

resource "thetalake_integration" "api" {
  name           = "Custom Theta Lake API Integration"
  theta_lake_api = {}
}
