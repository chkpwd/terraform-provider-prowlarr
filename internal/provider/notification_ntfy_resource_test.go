package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationNtfyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationNtfyResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_ntfy.test", "password", "key1"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_ntfy.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationNtfyResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_ntfy.test", "password", "key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_ntfy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationNtfyResourceConfig(name, key string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_ntfy" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"

		priority = 1
		server_url = "https://ntfy.sh"
		username = "User"
		password = "%s"
		topics = ["Topic1234","Topic4321"]
		field_tags = ["warning","skull"]
	}`, name, key)
}
