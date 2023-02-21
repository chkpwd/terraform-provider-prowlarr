package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

const applicationResourceName = "application"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ApplicationResource{}
	_ resource.ResourceWithImportState = &ApplicationResource{}
)

var (
	applicationStringFields   = []string{"prowlarrUrl", "baseUrl", "apiKey"}
	applicationIntSliceFields = []string{"syncCategories", "animeSyncCategories"}
)

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

// ApplicationResource defines the application implementation.
type ApplicationResource struct {
	client *prowlarr.APIClient
}

// Application describes the application data model.
type Application struct {
	SyncCategories      types.Set    `tfsdk:"sync_categories"`
	AnimeSyncCategories types.Set    `tfsdk:"anime_sync_categories"`
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	ConfigContract      types.String `tfsdk:"config_contract"`
	Implementation      types.String `tfsdk:"implementation"`
	SyncLevel           types.String `tfsdk:"sync_level"`
	ProwlarrURL         types.String `tfsdk:"prowlarr_url"`
	BaseURL             types.String `tfsdk:"base_url"`
	APIKey              types.String `tfsdk:"api_key"`
	ID                  types.Int64  `tfsdk:"id"`
}

func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + applicationResourceName
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Applications -->Generic Application resource. When possible use a specific resource instead.\nFor more information refer to [Application](https://wiki.servarr.com/prowlarr/settings#applications).",
		Attributes: map[string]schema.Attribute{
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Application configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Application implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name.",
				Required:            true,
			},
			"sync_level": schema.StringAttribute{
				MarkdownDescription: "Sync level.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("addOnly", "disabled", "fullSync"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Application ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"prowlarr_url": schema.StringAttribute{
				MarkdownDescription: "Prowlarr URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"sync_categories": schema.SetAttribute{
				MarkdownDescription: "Sync categories.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"anime_sync_categories": schema.SetAttribute{
				MarkdownDescription: "Anime sync categories.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var application *Application

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Application
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.CreateApplications(ctx).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, applicationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+applicationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Application

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var application *Application

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Application current value
	response, _, err := r.client.ApplicationApi.GetApplicationsById(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+applicationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Application

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var application *Application

	resp.Diagnostics.Append(req.Plan.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Application
	request := application.read(ctx)

	response, _, err := r.client.ApplicationApi.UpdateApplications(ctx, strconv.Itoa(int(request.GetId()))).ApplicationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, applicationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+applicationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Application

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var application *Application

	resp.Diagnostics.Append(req.State.Get(ctx, &application)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Application current value
	_, err := r.client.ApplicationApi.DeleteApplications(ctx, int32(application.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, applicationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+applicationResourceName+": "+strconv.Itoa(int(application.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+applicationResourceName+": "+req.ID)
}

func (n *Application) write(ctx context.Context, application *prowlarr.ApplicationResource) {
	n.ID = types.Int64Value(int64(application.GetId()))
	n.Name = types.StringValue(application.GetName())
	n.SyncLevel = types.StringValue(string(application.GetSyncLevel()))
	n.Implementation = types.StringValue(application.GetImplementation())
	n.ConfigContract = types.StringValue(application.GetConfigContract())
	n.Tags = types.SetValueMust(types.Int64Type, nil)
	n.SyncCategories = types.SetValueMust(types.Int64Type, nil)
	n.AnimeSyncCategories = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, application.Tags, n.Tags.Type(ctx), &n.Tags)
	n.writeFields(ctx, application.GetFields())
}

func (n *Application) writeFields(ctx context.Context, fields []*prowlarr.Field) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(applicationStringFields, f.GetName()) {
			helpers.WriteStringField(f, n)

			continue
		}

		if slices.Contains(applicationIntSliceFields, f.GetName()) {
			helpers.WriteIntSliceField(ctx, f, n)
		}
	}
}

func (n *Application) read(ctx context.Context) *prowlarr.ApplicationResource {
	tags := make([]*int32, len(n.Tags.Elements()))
	tfsdk.ValueAs(ctx, n.Tags, &tags)

	application := prowlarr.NewApplicationResource()
	application.SetSyncLevel(prowlarr.ApplicationSyncLevel(n.SyncLevel.ValueString()))
	application.SetId(int32(n.ID.ValueInt64()))
	application.SetName(n.Name.ValueString())
	application.SetImplementation(n.Implementation.ValueString())
	application.SetConfigContract(n.ConfigContract.ValueString())
	application.SetTags(tags)
	application.SetFields(n.readFields(ctx))

	return application
}

func (n *Application) readFields(ctx context.Context) []*prowlarr.Field {
	var output []*prowlarr.Field

	for _, s := range applicationStringFields {
		if field := helpers.ReadStringField(s, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range applicationIntSliceFields {
		if field := helpers.ReadIntSliceField(ctx, s, n); field != nil {
			output = append(output, field)
		}
	}

	return output
}