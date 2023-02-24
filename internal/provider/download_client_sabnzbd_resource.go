package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientSabnzbdResourceName   = "download_client_sabnzbd"
	downloadClientSabnzbdImplementation = "Sabnzbd"
	downloadClientSabnzbdConfigContract = "SabnzbdSettings"
	downloadClientSabnzbdProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientSabnzbdResource{}
	_ resource.ResourceWithImportState = &DownloadClientSabnzbdResource{}
)

func NewDownloadClientSabnzbdResource() resource.Resource {
	return &DownloadClientSabnzbdResource{}
}

// DownloadClientSabnzbdResource defines the download client implementation.
type DownloadClientSabnzbdResource struct {
	client *prowlarr.APIClient
}

// DownloadClientSabnzbd describes the download client data model.
type DownloadClientSabnzbd struct {
	Tags         types.Set    `tfsdk:"tags"`
	Categories   types.Set    `tfsdk:"categories"`
	Name         types.String `tfsdk:"name"`
	Host         types.String `tfsdk:"host"`
	URLBase      types.String `tfsdk:"url_base"`
	APIKey       types.String `tfsdk:"api_key"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Category     types.String `tfsdk:"category"`
	ItemPriority types.Int64  `tfsdk:"item_priority"`
	Priority     types.Int64  `tfsdk:"priority"`
	Port         types.Int64  `tfsdk:"port"`
	ID           types.Int64  `tfsdk:"id"`
	UseSsl       types.Bool   `tfsdk:"use_ssl"`
	Enable       types.Bool   `tfsdk:"enable"`
}

func (d DownloadClientSabnzbd) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:           d.Tags,
		Categories:     d.Categories,
		Name:           d.Name,
		Host:           d.Host,
		URLBase:        d.URLBase,
		APIKey:         d.APIKey,
		Username:       d.Username,
		Password:       d.Password,
		Category:       d.Category,
		ItemPriority:   d.ItemPriority,
		Priority:       d.Priority,
		Port:           d.Port,
		ID:             d.ID,
		UseSsl:         d.UseSsl,
		Enable:         d.Enable,
		Implementation: types.StringValue(downloadClientSabnzbdImplementation),
		ConfigContract: types.StringValue(downloadClientSabnzbdConfigContract),
		Protocol:       types.StringValue(downloadClientSabnzbdProtocol),
	}
}

func (d *DownloadClientSabnzbd) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Categories = client.Categories
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.APIKey = client.APIKey
	d.Username = client.Username
	d.Password = client.Password
	d.Category = client.Category
	d.ItemPriority = client.ItemPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
}

func (r *DownloadClientSabnzbdResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientSabnzbdResourceName
}

func (r *DownloadClientSabnzbdResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Sabnzbd resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [Sabnzbd](https://wiki.servarr.com/prowlarr/supported#sabnzbd).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"categories": schema.SetNestedAttribute{
				MarkdownDescription: "List of mapped categories.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: DownloadClientResource{}.getClientCategorySchema().Attributes,
				},
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"item_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Movie priority. `-100` Default, `-2` Paused, `-1` Low, `0` Normal, `1` High, `2` Force.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-100, -2, -1, 0, 1, 2),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientSabnzbdResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientSabnzbdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientSabnzbd

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientSabnzbd
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientSabnzbdResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientSabnzbdResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientSabnzbdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientSabnzbd

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientSabnzbd current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientSabnzbdResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientSabnzbdResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientSabnzbdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientSabnzbd

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientSabnzbd
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientSabnzbdResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientSabnzbdResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientSabnzbdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientSabnzbd

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientSabnzbd current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientSabnzbdResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientSabnzbdResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientSabnzbdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientSabnzbdResourceName+": "+req.ID)
}

func (d *DownloadClientSabnzbd) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{}
	genericDownloadClient.write(ctx, downloadClient)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientSabnzbd) read(ctx context.Context) *prowlarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx)
}