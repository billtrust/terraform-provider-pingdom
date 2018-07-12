variable "pingdom_user" {}
variable "pingdom_password" {}
variable "pingdom_api_key" {}
variable "pingdom_account_email" {}

provider "pingdom" {
  user          = "${var.pingdom_user}"
  password      = "${var.pingdom_password}"
  api_key       = "${var.pingdom_api_key}"
  account_email = "${var.pingdom_account_email}" # Main account email
}

resource "pingdom_user" "test_user" {
  name = "testUsername"
  paused = "NO"
}

resource "pingdom_user_contact_email" "test_user_email" {
  user_id = "${pingdom_user.test_user.id}"
  severity = "HIGH"
  address = "test@example.com"
}

resource "pingdom_user_contact_sms" "test_user_sms" {
  user_id = "${pingdom_user.test_user.id}"
  severity = "LOW"
  number = "55555555555"
  country_code = "1"
  phone_provider = "nexmo"
}