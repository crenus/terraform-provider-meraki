---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_organizations_snmp Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  OrganizationsSnmp resource for updating org snmp settings.
---

# meraki_organizations_snmp (Resource)

OrganizationsSnmp resource for updating org snmp settings.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network Id

### Optional

- `access` (String) The type of SNMP access. Can be one of 'none' (disabled), 'community' (V1/V2c), or 'users' (V3).
- `users` (Attributes Set) The list of SNMP users. Only relevant if 'access' is set to 'users'. (see [below for nested schema](#nestedatt--users))

### Read-Only

- `id` (String) Example identifier

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Optional:

- `passphrase` (String) The passphrase for the SNMP user.
- `username` (String) The username for the SNMP user

