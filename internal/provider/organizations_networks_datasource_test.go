package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsNetworksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsNetworksDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks"),
				),
			},

			// Create and Read Network
			{
				Config: testAccOrganizationsNetworksDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Read OrganizationsNetworks
			{
				Config: testAccOrganizationsNetworksDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.name", "test_acc_network"),

					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.tags.#", "2"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.product_types.#", "3"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.product_types.0", "appliance"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.product_types.1", "switch"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.product_types.2", "wireless"),
					resource.TestCheckResourceAttr("data.meraki_organizations_networks.test", "list.0.notes", "Additional description of the network"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsNetworksDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks"
 	api_enabled = true
 }
 `

const testAccOrganizationsNetworksDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1", "tag2"]
	name = "test_acc_network"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccOrganizationsNetworksDataSourceConfigRead = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

data "meraki_organizations_networks" "test" {
	depends_on = ["meraki_organization.test", meraki_network.test]
	organization_id = resource.meraki_organization.test.organization_id
}
`
