package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"

	"github.com/devopsarr/terraform-provider-prowlarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

const notificationResourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationResource{}
	_ resource.ResourceWithImportState = &NotificationResource{}
)

var (
	notificationBoolFields        = []string{"alwaysUpdate", "cleanLibrary", "directMessage", "notify", "requireEncryption", "sendSilently", "useSsl", "updateLibrary", "useEuEndpoint"}
	notificationStringFields      = []string{"accessToken", "accessTokenSecret", "apiKey", "aPIKey", "appToken", "arguments", "author", "authToken", "authUser", "avatar", "cc", "bcc", "botToken", "channel", "chatId", "consumerKey", "consumerSecret", "deviceNames", "displayTime", "expires", "from", "host", "icon", "instanceName", "mention", "password", "path", "refreshToken", "senderDomain", "senderId", "server", "signIn", "sound", "to", "token", "url", "userKey", "username", "webHookUrl", "serverUrl", "userName", "clickUrl", "mapFrom", "mapTo", "key", "event"}
	notificationIntFields         = []string{"port", "grabFields", "importFields", "priority", "retry", "expire", "method"}
	notificationStringSliceFields = []string{"recipients", "topics", "tags", "channelTags", "devices"}
	notificationIntSliceFields    = []string{"deviceIds"}
)

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the notification implementation.
type NotificationResource struct {
	client *prowlarr.APIClient
}

// Notification describes the notification data model.
type Notification struct {
	Tags                  types.Set    `tfsdk:"tags"`
	FieldTags             types.Set    `tfsdk:"field_tags"`
	ChannelTags           types.Set    `tfsdk:"channel_tags"`
	Topics                types.Set    `tfsdk:"topics"`
	DeviceIds             types.Set    `tfsdk:"device_ids"`
	Devices               types.Set    `tfsdk:"devices"`
	Recipients            types.Set    `tfsdk:"recipients"`
	DeviceNames           types.String `tfsdk:"device_names"`
	AccessToken           types.String `tfsdk:"access_token"`
	Host                  types.String `tfsdk:"host"`
	InstanceName          types.String `tfsdk:"instance_name"`
	Name                  types.String `tfsdk:"name"`
	Implementation        types.String `tfsdk:"implementation"`
	ConfigContract        types.String `tfsdk:"config_contract"`
	ClickURL              types.String `tfsdk:"click_url"`
	ConsumerSecret        types.String `tfsdk:"consumer_secret"`
	Path                  types.String `tfsdk:"path"`
	Arguments             types.String `tfsdk:"arguments"`
	ConsumerKey           types.String `tfsdk:"consumer_key"`
	ChatID                types.String `tfsdk:"chat_id"`
	From                  types.String `tfsdk:"from"`
	Cc                    types.String `tfsdk:"cc"`
	Icon                  types.String `tfsdk:"icon"`
	Password              types.String `tfsdk:"password"`
	Event                 types.String `tfsdk:"event"`
	Key                   types.String `tfsdk:"key"`
	DisplayTime           types.String `tfsdk:"display_time"`
	RefreshToken          types.String `tfsdk:"refresh_token"`
	WebHookURL            types.String `tfsdk:"web_hook_url"`
	Username              types.String `tfsdk:"username"`
	UserKey               types.String `tfsdk:"user_key"`
	Mention               types.String `tfsdk:"mention"`
	Avatar                types.String `tfsdk:"avatar"`
	URL                   types.String `tfsdk:"url"`
	Token                 types.String `tfsdk:"token"`
	To                    types.String `tfsdk:"to"`
	Sound                 types.String `tfsdk:"sound"`
	Bcc                   types.String `tfsdk:"bcc"`
	SignIn                types.String `tfsdk:"sign_in"`
	Server                types.String `tfsdk:"server"`
	SenderID              types.String `tfsdk:"sender_id"`
	BotToken              types.String `tfsdk:"bot_token"`
	SenderDomain          types.String `tfsdk:"sender_domain"`
	MapTo                 types.String `tfsdk:"map_to"`
	MapFrom               types.String `tfsdk:"map_from"`
	Channel               types.String `tfsdk:"channel"`
	Expires               types.String `tfsdk:"expires"`
	ServerURL             types.String `tfsdk:"server_url"`
	AccessTokenSecret     types.String `tfsdk:"access_token_secret"`
	APIKey                types.String `tfsdk:"api_key"`
	AppToken              types.String `tfsdk:"app_token"`
	Author                types.String `tfsdk:"author"`
	AuthToken             types.String `tfsdk:"auth_token"`
	AuthUser              types.String `tfsdk:"auth_user"`
	Priority              types.Int64  `tfsdk:"priority"`
	Port                  types.Int64  `tfsdk:"port"`
	Method                types.Int64  `tfsdk:"method"`
	Retry                 types.Int64  `tfsdk:"retry"`
	Expire                types.Int64  `tfsdk:"expire"`
	ID                    types.Int64  `tfsdk:"id"`
	ImportFields          types.Int64  `tfsdk:"import_fields"`
	GrabFields            types.Int64  `tfsdk:"grab_fields"`
	CleanLibrary          types.Bool   `tfsdk:"clean_library"`
	SendSilently          types.Bool   `tfsdk:"send_silently"`
	AlwaysUpdate          types.Bool   `tfsdk:"always_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	DirectMessage         types.Bool   `tfsdk:"direct_message"`
	RequireEncryption     types.Bool   `tfsdk:"require_encryption"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Notify                types.Bool   `tfsdk:"notify"`
	UseEuEndpoint         types.Bool   `tfsdk:"use_eu_endpoint"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
}

func (r *NotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationResourceName
}

func (r *NotificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification resource.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect).",
		Attributes: map[string]schema.Attribute{
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Notification configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Notification implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name.",
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
			"always_update": schema.BoolAttribute{
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
			},
			"clean_library": schema.BoolAttribute{
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
			},
			"direct_message": schema.BoolAttribute{
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
			},
			"require_encryption": schema.BoolAttribute{
				MarkdownDescription: "Require encryption flag.",
				Optional:            true,
				Computed:            true,
			},
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Add silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
			},
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
			"grab_fields": schema.Int64Attribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart, `10` CustomFormats, `11` CustomFormatScore.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
				},
			},
			"import_fields": schema.Int64Attribute{
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12),
				},
			},
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.", // TODO: add values in description
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2, 3, 4, 5, 7, 8),
				},
			},
			"retry": schema.Int64Attribute{
				MarkdownDescription: "Retry.",
				Optional:            true,
				Computed:            true,
			},
			"expire": schema.Int64Attribute{
				MarkdownDescription: "Expire.",
				Optional:            true,
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App token.",
				Optional:            true,
				Computed:            true,
			},
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Auth token.",
				Optional:            true,
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Optional:            true,
				Computed:            true,
			},
			"bcc": schema.StringAttribute{
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Optional:            true,
				Computed:            true,
			},
			"cc": schema.StringAttribute{
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_key": schema.StringAttribute{
				MarkdownDescription: "Consumer key.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_secret": schema.StringAttribute{
				MarkdownDescription: "Consumer secret.",
				Optional:            true,
				Computed:            true,
			},
			"device_names": schema.StringAttribute{
				MarkdownDescription: "Device names.",
				Optional:            true,
				Computed:            true,
			},
			"display_time": schema.StringAttribute{
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Optional:            true,
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
			},
			"mention": schema.StringAttribute{
				MarkdownDescription: "Mention.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Optional:            true,
				Computed:            true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
			},
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "server.",
				Optional:            true,
				Computed:            true,
			},
			"sign_in": schema.StringAttribute{
				MarkdownDescription: "Sign in.",
				Optional:            true,
				Computed:            true,
			},
			"sound": schema.StringAttribute{
				MarkdownDescription: "Sound.",
				Optional:            true,
				Computed:            true,
			},
			"to": schema.StringAttribute{
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "User key.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "Web hook url.",
				Optional:            true,
				Computed:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server url.",
				Optional:            true,
				Computed:            true,
			},
			"click_url": schema.StringAttribute{
				MarkdownDescription: "Click URL.",
				Optional:            true,
				Computed:            true,
			},
			"map_from": schema.StringAttribute{
				MarkdownDescription: "Map From.",
				Optional:            true,
				Computed:            true,
			},
			"map_to": schema.StringAttribute{
				MarkdownDescription: "Map To.",
				Optional:            true,
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Optional:            true,
				Computed:            true,
			},
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "Device IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "Channel tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"topics": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*prowlarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *prowlarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Notification current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *Notification) write(ctx context.Context, notification *prowlarr.NotificationResource) {
	n.OnHealthIssue = types.BoolValue(notification.GetOnHealthIssue())
	n.OnApplicationUpdate = types.BoolValue(notification.GetOnApplicationUpdate())
	n.IncludeHealthWarnings = types.BoolValue(notification.GetIncludeHealthWarnings())
	n.ID = types.Int64Value(int64(notification.GetId()))
	n.Name = types.StringValue(notification.GetName())
	n.Implementation = types.StringValue(notification.GetImplementation())
	n.ConfigContract = types.StringValue(notification.GetConfigContract())
	n.Tags = types.SetValueMust(types.Int64Type, nil)
	n.ChannelTags = types.SetValueMust(types.StringType, nil)
	n.DeviceIds = types.SetValueMust(types.Int64Type, nil)
	n.Topics = types.SetValueMust(types.StringType, nil)
	n.Devices = types.SetValueMust(types.StringType, nil)
	n.Recipients = types.SetValueMust(types.StringType, nil)
	n.FieldTags = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, notification.Tags, n.Tags.Type(ctx), &n.Tags)
	n.writeFields(ctx, notification.Fields)
}

func (n *Notification) writeFields(ctx context.Context, fields []*prowlarr.Field) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(notificationStringFields, f.GetName()) {
			tools.WriteStringField(f, n)

			continue
		}

		if slices.Contains(notificationBoolFields, f.GetName()) {
			tools.WriteBoolField(f, n)

			continue
		}

		if slices.Contains(notificationIntFields, f.GetName()) {
			tools.WriteIntField(f, n)

			continue
		}

		if slices.Contains(notificationStringSliceFields, f.GetName()) {
			tools.WriteStringSliceField(ctx, f, n)
		}

		if slices.Contains(notificationIntSliceFields, f.GetName()) {
			tools.WriteIntSliceField(ctx, f, n)
		}
	}
}

func (n *Notification) read(ctx context.Context) *prowlarr.NotificationResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := prowlarr.NewNotificationResource()
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetImplementation(n.Implementation.ValueString())
	notification.SetConfigContract(n.ConfigContract.ValueString())
	notification.SetTags(tags)
	notification.SetFields(n.readFields(ctx))

	return notification
}

func (n *Notification) readFields(ctx context.Context) []*prowlarr.Field {
	var output []*prowlarr.Field

	for _, b := range notificationBoolFields {
		if field := tools.ReadBoolField(b, n); field != nil {
			output = append(output, field)
		}
	}

	for _, i := range notificationIntFields {
		if field := tools.ReadIntField(i, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationStringFields {
		if field := tools.ReadStringField(s, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationStringSliceFields {
		if field := tools.ReadStringSliceField(ctx, s, n); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range notificationIntSliceFields {
		if field := tools.ReadIntSliceField(ctx, s, n); field != nil {
			output = append(output, field)
		}
	}

	return output
}