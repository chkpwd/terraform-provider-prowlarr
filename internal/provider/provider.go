package provider

import (
	"context"
	"os"

	"github.com/devopsarr/prowlarr-go/prowlarr"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &ProwlarrProvider{}

// ScaffoldingProvider defines the provider implementation.
type ProwlarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Prowlarr describes the provider data model.
type Prowlarr struct {
	APIKey        types.String `tfsdk:"api_key"`
	Authorization types.String `tfsdk:"authorization"`
	URL           types.String `tfsdk:"url"`
}

func (p *ProwlarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "prowlarr"
	resp.Version = p.version
}

func (p *ProwlarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Prowlarr provider is used to interact with any [Prowlarr](https://prowlarr.com/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Prowlarr authentication. Can be specified via the `PROWLARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"authorization": schema.StringAttribute{
				MarkdownDescription: "Token for token-based authentication with Prowlarr. This is an alternative to using an API key. Set this via the `PROWLARR_AUTHORIZATION` environment variable. One of `authorization` or `api_key` must be provided, but not both.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Prowlarr URL with protocol and port (e.g. `https://test.prowlarr.com:9696`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `PROWLARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *ProwlarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Prowlarr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("PROWLARR_URL")
	} else {
		url = data.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("PROWLARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	var authorization string
	if data.Authorization.IsNull() {
		authorization = os.Getenv("PROWLARR_AUTHORIZATION")
	} else {
		authorization = data.Authorization.ValueString()
	}

	if key == "" && authorization == "" {
		resp.Diagnostics.AddError(
			"Missing Authentication Credentials",
			"Both 'api_key' and 'authorization' are empty. You must provide either an API key or an authorization token for Prowlarr authentication.",
		)

		return
	}

	if key != "" && authorization != "" {
		resp.Diagnostics.AddError(
			"Conflicting Authentication Credentials",
			"Both 'api_key' and 'authorization' are provided. You must only provide one of these for Prowlarr authentication",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := prowlarr.NewConfiguration()

	if key != "" {
		config.AddDefaultHeader("X-Api-Key", key)
	}

	if authorization != "" {
		config.AddDefaultHeader("Authorization", authorization)
	}

	config.Servers[0].URL = url
	client := prowlarr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ProwlarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Applications
		NewSyncProfileResource,
		NewApplicationResource,
		NewApplicationLazyLibrarianResource,
		NewApplicationLidarrResource,
		NewApplicationMylarResource,
		NewApplicationRadarrResource,
		NewApplicationReadarrResource,
		NewApplicationSonarrResource,
		NewApplicationWhisparrResource,

		// Download Clients
		NewDownloadClientResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientFreeboxResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,

		// Indexer Proxies
		NewIndexerProxyResource,
		NewIndexerProxyFlaresolverrResource,
		NewIndexerProxyHTTPResource,
		NewIndexerProxySocks4Resource,
		NewIndexerProxySocks5Resource,

		// Indexer
		NewIndexerResource,

		// Notifications
		NewNotificationResource,
		NewNotificationAppriseResource,
		NewNotificationBoxcarResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationNtfyResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSignalResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationTelegramResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
	}
}

func (p *ProwlarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Applications
		NewSyncProfileDataSource,
		NewSyncProfilesDataSource,
		NewApplicationDataSource,
		NewApplicationsDataSource,

		// Download Clients
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,

		// Indexer Proxies
		NewIndexerProxyDataSource,
		NewIndexerProxiesDataSource,

		// Indexer
		NewIndexerDataSource,
		NewIndexersDataSource,
		NewIndexerSchemaDataSource,
		NewIndexerSchemasDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// System
		NewHostDataSource,
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
		NewTagDetailsDataSource,
		NewTagsDetailsDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ProwlarrProvider{
			version: version,
		}
	}
}
