---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_notification_email Resource - terraform-provider-prowlarr"
subcategory: "Notifications"
description: |-
  Notification Email resource.
  For more information refer to Notification https://wiki.servarr.com/prowlarr/settings#connect and Email https://wiki.servarr.com/prowlarr/supported#email.
---

# prowlarr_notification_email (Resource)

<!-- subcategory:Notifications -->Notification Email resource.
For more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Email](https://wiki.servarr.com/prowlarr/supported#email).

## Example Usage

```terraform
resource "prowlarr_notification_email" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  server = "http://email-server.net"
  port   = 587
  from   = "from_email@example.com"
  to     = ["user1@example.com", "user2@example.com"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `from` (String) From.
- `name` (String) NotificationEmail name.
- `server` (String) Server.
- `to` (Set of String) To.

### Optional

- `bcc` (Set of String) Bcc.
- `cc` (Set of String) Cc.
- `include_health_warnings` (Boolean) Include health warnings.
- `include_manual_grabs` (Boolean) Include manual grab flag.
- `on_application_update` (Boolean) On application update flag.
- `on_grab` (Boolean) On release grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_health_restored` (Boolean) On health restored flag.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `require_encryption` (Boolean) Require encryption flag.
- `tags` (Set of Number) List of associated tags.
- `username` (String) Username.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_notification_email.example 1
```
