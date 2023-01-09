---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_download_client Resource - terraform-provider-prowlarr"
subcategory: "Download Clients"
description: |-
  Download Client resource.
  For more information refer to Download Client https://wiki.servarr.com/prowlarr/settings#download-clients.
---

# prowlarr_download_client (Resource)

<!-- subcategory:Download Clients -->Download Client resource.
For more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients).

## Example Usage

```terraform
resource "prowlarr_download_client" "example" {
  enable          = true
  priority        = 1
  name            = "Example"
  implementation  = "Transmission"
  protocol        = "torrent"
  config_contract = "TransmissionSettings"
  host            = "transmission"
  url_base        = "/transmission/"
  port            = 9091
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config_contract` (String) DownloadClient configuration template.
- `implementation` (String) DownloadClient implementation name.
- `name` (String) Download Client name.
- `protocol` (String) Protocol. Valid values are 'usenet' and 'torrent'.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `add_stopped` (Boolean) Add stopped flag.
- `additional_tags` (Set of Number) Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.
- `api_key` (String) API key.
- `categories` (Attributes Set) List of mapped categories. (see [below for nested schema](#nestedatt--categories))
- `category` (String) Category.
- `destination` (String) Destination.
- `directory` (String) Directory.
- `enable` (Boolean) Enable flag.
- `field_tags` (Set of String) Field tags.
- `first_and_last` (Boolean) First and last flag.
- `host` (String) host.
- `initial_state` (Number) Initial state. `0` Start, `1` ForceStart, `2` Pause.
- `intial_state` (Number) Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.
- `item_priority` (Number) Priority. `0` Last, `1` First.
- `magnet_file_extension` (String) Magnet file extension.
- `nzb_folder` (String) NZB folder.
- `password` (String) Password.
- `port` (Number) Port.
- `post_im_tags` (Set of String) Post import tags.
- `priority` (Number) Priority.
- `read_only` (Boolean) Read only flag.
- `rpc_path` (String) RPC path.
- `save_magnet_files` (Boolean) Save magnet files flag.
- `secret_token` (String) Secret token.
- `sequential_order` (Boolean) Sequential order flag.
- `start_on_add` (Boolean) Start on add flag.
- `strm_folder` (String) STRM folder.
- `tags` (Set of Number) List of associated tags.
- `torrent_folder` (String) Torrent folder.
- `tv_imported_category` (String) TV imported category.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.
- `watch_folder` (Boolean) Watch folder flag.

### Read-Only

- `id` (Number) Download Client ID.

<a id="nestedatt--categories"></a>
### Nested Schema for `categories`

Optional:

- `categories` (Set of Number) List of categories.
- `name` (String) Name of client category.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_download_client.example 1
```