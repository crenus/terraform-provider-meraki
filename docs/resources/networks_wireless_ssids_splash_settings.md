---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_wireless_ssids_splash_settings Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  NetworksWirelessSsidsSplashSettings
---

# meraki_networks_wireless_ssids_splash_settings (Resource)

NetworksWirelessSsidsSplashSettings



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network Id
- `number` (String) SsIds Number

### Optional

- `allow_simultaneous_logins` (Boolean) Whether or not to allow simultaneous logins from different devices.
- `billing` (Attributes) Details associated with billing splash. (see [below for nested schema](#nestedatt--billing))
- `block_all_traffic_before_sign_on` (Boolean) How restricted allowing traffic should be. If true, all traffic types are blocked until the splash page is acknowledged. If false, all non-HTTP traffic is allowed before the splash page is acknowledged.
- `controller_disconnection_behavior` (String) How login attempts should be handled when the controller is unreachable. Can be either 'open', 'restricted', or 'default'.
- `guest_sponsorship` (Attributes) Details associated with guest sponsored splash. (see [below for nested schema](#nestedatt--guest_sponsorship))
- `redirect_url` (String) The custom redirect URL where the users will go after the splash page.
- `sentry_enrollment` (Attributes) Systems Manager sentry enrollment splash settings. (see [below for nested schema](#nestedatt--sentry_enrollment))
- `splash_image` (Attributes) The image used in the splash page. (see [below for nested schema](#nestedatt--splash_image))
- `splash_logo` (Attributes) The logo used in the splash page. (see [below for nested schema](#nestedatt--splash_logo))
- `splash_prepaid_front` (Attributes) The prepaid front image used in the splash page. (see [below for nested schema](#nestedatt--splash_prepaid_front))
- `splash_timeout` (Number) Splash timeout in minutes. This will determine how often users will see the splash page.
- `splash_url` (String) The custom splash URL of the click-through splash page. Note that the URL can be configured without necessarily being used. In order to enable the custom URL, see 'useSplashUrl'
- `use_redirect_url` (Boolean) Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.
- `use_splash_url` (Boolean) Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.
- `welcome_message` (String) The welcome message for the users on the splash page.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--billing"></a>
### Nested Schema for `billing`

Optional:

- `free_access` (Attributes) Details associated with a free access plan with limits. (see [below for nested schema](#nestedatt--billing--free_access))
- `prepaid_access_fast_login_enabled` (Boolean) Whether or not billing uses the fast login prepaid access option.
- `reply_to_email_address` (String) The email address that receives replies from clients.

<a id="nestedatt--billing--free_access"></a>
### Nested Schema for `billing.free_access`

Optional:

- `duration_in_minutes` (Number) How long a device can use a network for free..
- `enabled` (Boolean) Whether or not free access is enabled.



<a id="nestedatt--guest_sponsorship"></a>
### Nested Schema for `guest_sponsorship`

Optional:

- `duration_in_minutes` (Number) Duration in minutes of sponsored guest authorization. Must be between 1 and 60480 (6 weeks).
- `guest_can_request_time_frame` (Boolean) Whether or not guests can specify how much time they are requesting.


<a id="nestedatt--sentry_enrollment"></a>
### Nested Schema for `sentry_enrollment`

Optional:

- `enforced_systems` (Set of String) The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.
- `strength` (String) The strength of the enforcement of selected system types. Must be one of: 'focused', 'click-through', and 'strict'.
- `systems_manager_network` (Attributes) Systems Manager network targeted for sentry enrollment.. (see [below for nested schema](#nestedatt--sentry_enrollment--systems_manager_network))

<a id="nestedatt--sentry_enrollment--systems_manager_network"></a>
### Nested Schema for `sentry_enrollment.systems_manager_network`

Optional:

- `id` (String) The network ID of the Systems Manager network.



<a id="nestedatt--splash_image"></a>
### Nested Schema for `splash_image`

Optional:

- `extension` (String) The extension of the image file.
- `image` (Attributes) Properties for setting a new image. (see [below for nested schema](#nestedatt--splash_image--image))
- `md5` (String) The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.

<a id="nestedatt--splash_image--image"></a>
### Nested Schema for `splash_image.image`

Optional:

- `contents` (String) The file contents (a base 64 encoded string) of your new image.
- `format` (String) The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.



<a id="nestedatt--splash_logo"></a>
### Nested Schema for `splash_logo`

Optional:

- `extension` (String) The extension of the logo file.
- `image` (Attributes) Properties for setting a new image. (see [below for nested schema](#nestedatt--splash_logo--image))
- `md5` (String) The MD5 value of the logo file. Setting this to null will remove the logo from the splash page.

<a id="nestedatt--splash_logo--image"></a>
### Nested Schema for `splash_logo.image`

Optional:

- `contents` (String) The file contents (a base 64 encoded string) of your new logo.
- `format` (String) The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.



<a id="nestedatt--splash_prepaid_front"></a>
### Nested Schema for `splash_prepaid_front`

Optional:

- `extension` (String) The extension of the prepaid front image file.
- `image` (Attributes) Properties for setting a new image. (see [below for nested schema](#nestedatt--splash_prepaid_front--image))
- `md5` (String) The MD5 value of the prepaid front image file. Setting this to null will remove the prepaid front from the splash page.

<a id="nestedatt--splash_prepaid_front--image"></a>
### Nested Schema for `splash_prepaid_front.image`

Optional:

- `contents` (String) The file contents (a base 64 encoded string) of your new prepaid.
- `format` (String) The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.