package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdaptivePolicyAclsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create test organization
			{
				Config: testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_admin_adaptive_policy_acls"),
				),
			},

			// Create an adaptive policy acl
			{
				Config: testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateAcl,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "name", "Block sensitive web traffic"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "description", "Blocks sensitive web traffic"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "ip_version", "ipv6"),

					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.src_port", "1,33"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.dst_port", "22-30"),
				),
			},

			// Read an adaptive policy acl
			{
				Config: testAccOrganizationsAdaptivePolicyAclsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.0.rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.0.rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.0.rules.0.src_port", "1,33"),
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.0.rules.0.dst_port", "22-30"),
				),
			},
		},
	})
}

const testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateOrg = `
resource "meraki_organization" "test" {
	name = "test_meraki_organizations_admin_adaptive_policy_acls"
	api_enabled = true
}
`

// Create an adaptive policy acl
const testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateAcl = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_adaptive_policy_acl" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	name = "Block sensitive web traffic"
	description = "Blocks sensitive web traffic"
	ip_version   = "ipv6"
	rules = [      
		{
			"policy": "deny",
			"protocol": "tcp",
			"src_port": "1,33",
			"dst_port": "22-30"
		}
	]  
  }
`

const testAccOrganizationsAdaptivePolicyAclsDataSourceConfig = `
resource "meraki_organization" "test" {}

data "meraki_organizations_adaptive_policy_acls" "test" {
    organization_id = resource.meraki_organization.test.organization_id
}
`
