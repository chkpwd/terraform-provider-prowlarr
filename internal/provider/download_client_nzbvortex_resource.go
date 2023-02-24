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
	downloadClientNzbvortexResourceName   = "download_client_nzbvortex"
	downloadClientNzbvortexImplementation = "NzbVortex"
	downloadClientNzbvortexConfigContract = "NzbVortexSettings"
	downloadClientNzbvortexProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientNzbvortexResource{}
	_ resource.ResourceWithImportState = &DownloadClientNzbvortexResource{}
)

func NewDownloadClientNzbvortexResource() resource.Resource {
	return &DownloadClientNzbvortexResource{}
}

// DownloadClientNzbvortexResource defines the download client implementation.
type DownloadClientNzbvortexResource struct {
	client *prowlarr.APIClient
}

// DownloadClientNzbvortex describes the download client data model.
type DownloadClientNzbvortex struct {
	Tags         types.Set    `tfsdk:"tags"`
	Categories   types.Set    `tfsdk:"categories"`
	Name         types.String `tfsdk:"name"`
	Host         types.String `tfsdk:"host"`
	URLBase      types.String `tfsdk:"url_base"`
	APIKey       types.String `tfsdk:"api_key"`
	Category     types.String `tfsdk:"category"`
	ItemPriority types.Int64  `tfsdk:"item_priority"`
	Priority     types.Int64  `tfsdk:"priority"`
	Port         types.Int64  `tfsdk:"port"`
	ID           types.Int64  `tfsdk:"id"`
	UseSsl       types.Bool   `tfsdk:"use_ssl"`
	Enable       types.Bool   `tfsdk:"enable"`
}

func (d DownloadClientNzbvortex) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:           d.Tags,
		Categories:     d.Categories,
		Name:           d.Name,
		Host:           d.Host,
		URLBase:        d.URLBase,
		APIKey:         d.APIKey,
		Category:       d.Category,
		ItemPriority:   d.ItemPriority,
		Priority:       d.Priority,
		Port:           d.Port,
		ID:             d.ID,
		UseSsl:         d.UseSsl,
		Enable:         d.Enable,
		Implementation: types.StringValue(downloadClientNzbvortexImplementation),
		ConfigContract: types.StringValue(downloadClientNzbvortexConfigContract),
		Protocol:       types.StringValue(downloadClientNzbvortexProtocol),
	}
}

func (d *DownloadClientNzbvortex) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Categories = client.Categories
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.APIKey = client.APIKey
	d.Category = client.Category
	d.ItemPriority = client.ItemPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
}

func (r *DownloadClientNzbvortexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientNzbvortexResourceName
}

func (r *DownloadClientNzbvortexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Nzbvortex resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [Nzbvortex](https://wiki.servarr.com/prowlarr/supported#nzbvortex).",
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
				MarkdownDescription: "Recent Movie priority. `-1` Low, `0` Normal, `1` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-1, 0, 1),
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
				Required:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientNzbvortexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientNzbvortexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientNzbvortex
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientNzbvortex

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientNzbvortex current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientNzbvortex
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientNzbvortex current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientNzbvortexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientNzbvortexResourceName+": "+req.ID)
}

func (d *DownloadClientNzbvortex) write(ctx context.Context, downloadClient *prowlarr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{}
	genericDownloadClient.write(ctx, downloadClient)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientNzbvortex) read(ctx context.Context) *prowlarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx)
}