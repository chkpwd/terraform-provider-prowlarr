package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationPushoverResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationPushoverResourceConfig("error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationPushoverResourceConfig("resourcePushoverTest", 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_pushover.test", "priority", "0"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_pushover.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationPushoverResourceConfig("error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationPushoverResourceConfig("resourcePushoverTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_pushover.test", "priority", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_pushover.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationPushoverResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_pushover" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		api_key = "Key"
		user_key = "Test"
		priority = %d
	}`, name, priority)
}
