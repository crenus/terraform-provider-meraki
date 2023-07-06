---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_organizations_license Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  Organizations License Updates a license
---

# meraki_organizations_license (Resource)

Organizations License Updates a license



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `device_serial` (String) The serial number of the device to assign this license to. Set this to null to unassign the license. If a different license is already active on the device, this parameter will control queueing/dequeuing this license.
- `license_id` (String) License ID
- `organization_id` (String) Organization ID

### Optional

- `activation_date` (String) The date the license started burning.
- `claim_date` (String) The date the license was claimed into the organization.
- `duration_in_days` (Number) The duration of the individual license.
- `expiration_date` (String) The date the license will expire.
- `head_license_id` (String) The id of the head license this license is queued behind. If there is no head license, it returns nil.
- `license_key` (String) License Key.
- `license_type` (String) License Type.
- `network_id` (String) ID of the network the license is assigned to.
- `order_number` (String) Order Number.
- `permanently_queued_licenses` (Attributes) DEPRECATED List of permanently queued licenses attached to the license. Instead, use /organizations/{organizationId}/licenses?deviceSerial= to retrieved queued licenses for a given device. (see [below for nested schema](#nestedatt--permanently_queued_licenses))
- `seat_count` (Number) The number of seats of the license. Only applicable to SM licenses.
- `state` (String) The state of the license. All queued licenses have a status of `recentlyQueued`.
- `total_duration_in_days` (Number) The duration of the license plus all permanently queued licenses associated with it.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--permanently_queued_licenses"></a>
### Nested Schema for `permanently_queued_licenses`

Optional:

- `duration_in_days` (Number) The duration of the individual license.
- `id` (String) Permanently queued license ID.
- `license_key` (String) License key.
- `license_type` (String) License type.
- `order_number` (String) Order number.

