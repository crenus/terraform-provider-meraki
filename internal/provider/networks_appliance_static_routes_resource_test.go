package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceStaticRoutesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_static_routes"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read Networks Appliance Static Routes.
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateNetworksApplianceStaticRoutes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "name", "My route"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.129.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.128.1"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "name", "My route"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.129.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "enable", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_address", "22:33:44:55:66:77"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_ip_address", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_name", "Some client name"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.comment", "A reserved IP range"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.start", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.end", "192.168.128.2"),
				),
			},
		},
	})
}

const testAccNetworksApplianceStaticRoutesResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
    name = "test_acc_meraki_organizations_networks_appliance_static_routes"
    api_enabled = true
}
`
const testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    organization_id = resource.meraki_organization.test.organization_id
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`

const testAccNetworksApplianceStaticRoutesResourceConfigCreateNetworksApplianceStaticRoutes = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id  
    name = "My route"
    subnet = "192.168.129.0/24"
    gateway_ip = "192.168.128.1"
	reserved_ip_ranges = []
	
}
`

const testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id    
	name = "My route"
    subnet = "192.168.129.0/24"
	fixed_ip_assignments_mac_address = "22:33:44:55:66:77"
	fixed_ip_assignments_mac_ip_address = "192.168.128.1"
	fixed_ip_assignments_mac_name = "Some client name"   
	reserved_ip_ranges = [
        {
            start = "192.168.128.1"
            end = "192.168.128.2"
            comment = "A reserved IP range"
        }
    ]
}
`
