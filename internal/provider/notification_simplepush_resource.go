package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationSimplepushResourceName   = "notification_simplepush"
	notificationSimplepushImplementation = "Simplepush"
	notificationSimplepushConfigContract = "SimplepushSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSimplepushResource{}
	_ resource.ResourceWithImportState = &NotificationSimplepushResource{}
)

func NewNotificationSimplepushResource() resource.Resource {
	return &NotificationSimplepushResource{}
}

// NotificationSimplepushResource defines the notification implementation.
type NotificationSimplepushResource struct {
	client *prowlarr.APIClient
}

// NotificationSimplepush describes the notification data model.
type NotificationSimplepush struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Event                 types.String `tfsdk:"event"`
	Name                  types.String `tfsdk:"name"`
	Key                   types.String `tfsdk:"key"`
	ID                    types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
}

func (n NotificationSimplepush) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Event:                 n.Event,
		Key:                   n.Key,
		Name:                  n.Name,
		ID:                    n.ID,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		ConfigContract:        types.StringValue(notificationSimplepushConfigContract),
		Implementation:        types.StringValue(notificationSimplepushImplementation),
	}
}

func (n *NotificationSimplepush) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Event = notification.Event
	n.Key = notification.Key
	n.Name = notification.Name
	n.ID = notification.ID
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
}

func (r *NotificationSimplepushResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSimplepushResourceName
}

func (r *NotificationSimplepushResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Simplepush resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Simplepush](https://wiki.servarr.com/prowlarr/supported#simplepush).",
		Attributes: map[string]schema.Attribute{
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSimplepush name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationSimplepushResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationSimplepushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSimplepush
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSimplepush current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSimplepush
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSimplepush current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSimplepushResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSimplepushResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSimplepushResourceName+": "+req.ID)
}

func (n *NotificationSimplepush) write(ctx context.Context, notification *prowlarr.NotificationResource) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification)
	n.fromNotification(genericNotification)
}

func (n *NotificationSimplepush) read(ctx context.Context) *prowlarr.NotificationResource {
	return n.toNotification().read(ctx)
}
