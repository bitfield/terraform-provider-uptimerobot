variable "uptimerobot_api_key" {}

provider "uptimerobot" {
  api_key = var.uptimerobot_api_key
}

resource "uptimerobot_monitor" "test-monitor" {
  friendly_name = "My test monitor"
  url           = "http://bitfieldconsulting.com/"
  type          = "HTTP"
  alert_contact = ["2416450"]
}
