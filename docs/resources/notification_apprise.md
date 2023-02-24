---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_notification_apprise Resource - terraform-provider-prowlarr"
subcategory: "Notifications"
description: |-
  Notification Apprise resource.
  For more information refer to Notification https://wiki.servarr.com/prowlarr/settings#connect and Apprise https://wiki.servarr.com/prowlarr/supported#apprise.
---

# prowlarr_notification_apprise (Resource)

<!-- subcategory:Notifications -->Notification Apprise resource.
For more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Apprise](https://wiki.servarr.com/prowlarr/supported#apprise).

## Example Usage

```terraform
resource "prowlarr_notification_apprise" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  base_url          = "http://localhost:8000"
  configuration_key = "ConfigKey"
  auth_username     = "User"
  auth_password     = "Pass"
  field_tags        = ["test", "test1"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) NotificationApprise name.

### Optional

- `auth_password` (String, Sensitive) AuthPassword.
- `auth_username` (String) AuthUsername.
- `base_url` (String) Base URL.
- `configuration_key` (String, Sensitive) ConfigurationKey.
- `field_tags` (Set of String) Tags and emojis.
- `include_health_warnings` (Boolean) Include health warnings.
- `on_application_update` (Boolean) On application update flag.
- `on_health_issue` (Boolean) On health issue flag.
- `stateless_urls` (String) Comma separated stateless URLs.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_notification_apprise.example 1
```