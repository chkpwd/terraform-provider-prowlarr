package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationTelegramResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationTelegramResourceConfig("Error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_telegram.test", "chat_id", "chat01"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_telegram.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationTelegramResourceConfig("Error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_telegram.test", "chat_id", "chat02"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_telegram.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTelegramResourceConfig(name, chat string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_telegram" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		chat_id = "%s"
		bot_token = "Token"
	}`, name, chat)
}
