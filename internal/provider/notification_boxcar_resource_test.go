package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationBoxcarResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationBoxcarResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationBoxcarResourceConfig("resourceBoxcarTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_boxcar.test", "token", "token123"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_boxcar.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationBoxcarResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationBoxcarResourceConfig("resourceBoxcarTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_boxcar.test", "token", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_boxcar.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationBoxcarResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_boxcar" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		token = "%s"
	}`, name, token)
}
