package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccAdministeredIdentitiesMeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Read AdministeredIdentitiesMe
			{
				Config: testAccAdministeredIdentitiesMeDataSourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(

					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "name", ""),
					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "email", ""),
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "id", "example-id"),

					resource.TestCheckResourceAttrWith(
						"data.meraki_administered_identities_me.test", "last_used_dashboard_at", func(value string) error {

							re := regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)
							if re.MatchString(value) != true {
								return fmt.Errorf(fmt.Sprintf("received tiemstring does not match RFC3339 format: %s", value))
							}

							return nil
						}),
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "authentication_api_key_created", "true"),
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "authentication_mode", "email"),
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "authentication_two_factor_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "authentication_saml_enabled", "false"),
				),
			},
		},
	})
}

const testAccAdministeredIdentitiesMeDataSourceConfigCreate = `
data "meraki_administered_identities_me" "test" {
}
`
