package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const notificationResourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationResource{}
	_ resource.ResourceWithImportState = &NotificationResource{}
)

var notificationFields = helpers.Fields{
	Bools:                  []string{"alwaysUpdate", "cleanLibrary", "directMessage", "notify", "requireEncryption", "sendSilently", "useSsl", "updateLibrary", "useEuEndpoint"},
	Strings:                []string{"authPassword", "authUsername", "statelessUrls", "configurationKey", "baseUrl", "accessToken", "accessTokenSecret", "apiKey", "aPIKey", "appToken", "arguments", "author", "authToken", "authUser", "avatar", "botToken", "channel", "chatId", "consumerKey", "consumerSecret", "deviceNames", "expires", "from", "host", "icon", "instanceName", "mention", "password", "path", "refreshToken", "senderDomain", "senderId", "server", "signIn", "sound", "token", "url", "userKey", "username", "webHookUrl", "serverUrl", "userName", "clickUrl", "mapFrom", "mapTo", "key", "event", "topicId", "senderNumber", "receiverId"},
	Ints:                   []string{"displayTime", "port", "itemPriority", "retry", "expire", "method", "notificationType"},
	IntsExceptions:         []string{"priority"},
	StringSlices:           []string{"recipients", "to", "cC", "bcc", "topics", "fieldTags", "channelTags", "deviceIds", "devices"},
	StringSlicesExceptions: []string{"tags"},
	IntSlices:              []string{"grabFields"},
}

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
	GrabFields            types.Set    `tfsdk:"grab_fields"`
	DeviceIds             types.Set    `tfsdk:"device_ids"`
	Devices               types.Set    `tfsdk:"devices"`
	To                    types.Set    `tfsdk:"to"`
	Cc                    types.Set    `tfsdk:"cc"`
	Bcc                   types.Set    `tfsdk:"bcc"`
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
	TopicID               types.String `tfsdk:"topic_id"`
	From                  types.String `tfsdk:"from"`
	Icon                  types.String `tfsdk:"icon"`
	Password              types.String `tfsdk:"password"`
	Event                 types.String `tfsdk:"event"`
	Key                   types.String `tfsdk:"key"`
	RefreshToken          types.String `tfsdk:"refresh_token"`
	WebHookURL            types.String `tfsdk:"web_hook_url"`
	Username              types.String `tfsdk:"username"`
	UserKey               types.String `tfsdk:"user_key"`
	Mention               types.String `tfsdk:"mention"`
	Avatar                types.String `tfsdk:"avatar"`
	URL                   types.String `tfsdk:"url"`
	Token                 types.String `tfsdk:"token"`
	Sound                 types.String `tfsdk:"sound"`
	SignIn                types.String `tfsdk:"sign_in"`
	Server                types.String `tfsdk:"server"`
	SenderID              types.String `tfsdk:"sender_id"`
	SenderNumber          types.String `tfsdk:"sender_number"`
	ReceiverID            types.String `tfsdk:"receiver_id"`
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
	ConfigurationKey      types.String `tfsdk:"configuration_key"`
	StatelessURLs         types.String `tfsdk:"stateless_urls"`
	BaseURL               types.String `tfsdk:"base_url"`
	AuthUsername          types.String `tfsdk:"auth_username"`
	AuthPassword          types.String `tfsdk:"auth_password"`
	DisplayTime           types.Int64  `tfsdk:"display_time"`
	ItemPriority          types.Int64  `tfsdk:"priority"`
	Port                  types.Int64  `tfsdk:"port"`
	Method                types.Int64  `tfsdk:"method"`
	Retry                 types.Int64  `tfsdk:"retry"`
	Expire                types.Int64  `tfsdk:"expire"`
	NotificationType      types.Int64  `tfsdk:"notification_type"`
	ID                    types.Int64  `tfsdk:"id"`
	CleanLibrary          types.Bool   `tfsdk:"clean_library"`
	SendSilently          types.Bool   `tfsdk:"send_silently"`
	AlwaysUpdate          types.Bool   `tfsdk:"always_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
	DirectMessage         types.Bool   `tfsdk:"direct_message"`
	RequireEncryption     types.Bool   `tfsdk:"require_encryption"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Notify                types.Bool   `tfsdk:"notify"`
	UseEuEndpoint         types.Bool   `tfsdk:"use_eu_endpoint"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	IncludeManualGrabs    types.Bool   `tfsdk:"include_manual_grabs"`
}

func (n Notification) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":                    types.SetType{}.WithElementType(types.Int64Type),
			"grab_fields":             types.SetType{}.WithElementType(types.Int64Type),
			"device_ids":              types.SetType{}.WithElementType(types.Int64Type),
			"field_tags":              types.SetType{}.WithElementType(types.StringType),
			"recipients":              types.SetType{}.WithElementType(types.StringType),
			"devices":                 types.SetType{}.WithElementType(types.StringType),
			"to":                      types.SetType{}.WithElementType(types.StringType),
			"cc":                      types.SetType{}.WithElementType(types.StringType),
			"bcc":                     types.SetType{}.WithElementType(types.StringType),
			"channel_tags":            types.SetType{}.WithElementType(types.StringType),
			"topics":                  types.SetType{}.WithElementType(types.StringType),
			"device_names":            types.StringType,
			"access_token":            types.StringType,
			"host":                    types.StringType,
			"instance_name":           types.StringType,
			"name":                    types.StringType,
			"implementation":          types.StringType,
			"config_contract":         types.StringType,
			"click_url":               types.StringType,
			"consumer_secret":         types.StringType,
			"path":                    types.StringType,
			"arguments":               types.StringType,
			"consumer_key":            types.StringType,
			"chat_id":                 types.StringType,
			"topic_id":                types.StringType,
			"from":                    types.StringType,
			"icon":                    types.StringType,
			"password":                types.StringType,
			"event":                   types.StringType,
			"key":                     types.StringType,
			"refresh_token":           types.StringType,
			"web_hook_url":            types.StringType,
			"username":                types.StringType,
			"user_key":                types.StringType,
			"mention":                 types.StringType,
			"avatar":                  types.StringType,
			"url":                     types.StringType,
			"token":                   types.StringType,
			"sound":                   types.StringType,
			"sign_in":                 types.StringType,
			"server":                  types.StringType,
			"sender_id":               types.StringType,
			"sender_number":           types.StringType,
			"receiver_id":             types.StringType,
			"bot_token":               types.StringType,
			"sender_domain":           types.StringType,
			"map_to":                  types.StringType,
			"map_from":                types.StringType,
			"channel":                 types.StringType,
			"expires":                 types.StringType,
			"server_url":              types.StringType,
			"access_token_secret":     types.StringType,
			"api_key":                 types.StringType,
			"app_token":               types.StringType,
			"author":                  types.StringType,
			"auth_token":              types.StringType,
			"auth_user":               types.StringType,
			"configuration_key":       types.StringType,
			"stateless_urls":          types.StringType,
			"base_url":                types.StringType,
			"auth_username":           types.StringType,
			"auth_password":           types.StringType,
			"display_time":            types.Int64Type,
			"priority":                types.Int64Type,
			"port":                    types.Int64Type,
			"method":                  types.Int64Type,
			"retry":                   types.Int64Type,
			"expire":                  types.Int64Type,
			"notification_type":       types.Int64Type,
			"id":                      types.Int64Type,
			"clean_library":           types.BoolType,
			"send_silently":           types.BoolType,
			"always_update":           types.BoolType,
			"on_health_issue":         types.BoolType,
			"on_health_restored":      types.BoolType,
			"direct_message":          types.BoolType,
			"require_encryption":      types.BoolType,
			"use_ssl":                 types.BoolType,
			"notify":                  types.BoolType,
			"use_eu_endpoint":         types.BoolType,
			"update_library":          types.BoolType,
			"include_health_warnings": types.BoolType,
			"on_application_update":   types.BoolType,
			"on_grab":                 types.BoolType,
			"include_manual_grabs":    types.BoolType,
		})
}

func (r *NotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationResourceName
}

func (r *NotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Generic Notification resource. When possible use a specific resource instead.\nFor more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect).",
		Attributes: map[string]schema.Attribute{
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On release grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_manual_grabs": schema.BoolAttribute{
				MarkdownDescription: "Include manual grab flag.",
				Optional:            true,
				Computed:            true,
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
			"display_time": schema.Int64Attribute{
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
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
			"notification_type": schema.Int64Attribute{
				MarkdownDescription: "Notification type. `0` Info, `1` Success, `2` Warning, `3` Failure.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3),
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"stateless_urls": schema.StringAttribute{
				MarkdownDescription: "Comma separated stateless URLs.",
				Optional:            true,
				Computed:            true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "Auth username.",
				Optional:            true,
				Computed:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "Auth password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"configuration_key": schema.StringAttribute{
				MarkdownDescription: "Configuration key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Optional:            true,
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
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
			"topic_id": schema.StringAttribute{
				MarkdownDescription: "Topic ID.",
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
				Sensitive:           true,
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
			"sender_number": schema.StringAttribute{
				MarkdownDescription: "Sender Number.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"receiver_id": schema.StringAttribute{
				MarkdownDescription: "Receiver ID.",
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
			"grab_fields": schema.SetAttribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
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
			"to": schema.SetAttribute{
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cc": schema.SetAttribute{
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"bcc": schema.SetAttribute{
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response, &resp.Diagnostics)
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
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response, &resp.Diagnostics)
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
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationResourceName+": "+req.ID)
}

func (n *Notification) write(ctx context.Context, notification *prowlarr.NotificationResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	n.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	diags.Append(localDiag...)

	n.OnHealthIssue = types.BoolValue(notification.GetOnHealthIssue())
	n.OnHealthRestored = types.BoolValue(notification.GetOnHealthRestored())
	n.OnApplicationUpdate = types.BoolValue(notification.GetOnApplicationUpdate())
	n.OnGrab = types.BoolValue(notification.GetOnGrab())
	n.IncludeManualGrabs = types.BoolValue(notification.GetIncludeManualGrabs())
	n.IncludeHealthWarnings = types.BoolValue(notification.GetIncludeHealthWarnings())
	n.ID = types.Int64Value(int64(notification.GetId()))
	n.Name = types.StringValue(notification.GetName())
	n.Implementation = types.StringValue(notification.GetImplementation())
	n.ConfigContract = types.StringValue(notification.GetConfigContract())
	n.GrabFields = types.SetValueMust(types.Int64Type, nil)
	n.ChannelTags = types.SetValueMust(types.StringType, nil)
	n.DeviceIds = types.SetValueMust(types.Int64Type, nil)
	n.Topics = types.SetValueMust(types.StringType, nil)
	n.Devices = types.SetValueMust(types.StringType, nil)
	n.Recipients = types.SetValueMust(types.StringType, nil)
	n.FieldTags = types.SetValueMust(types.StringType, nil)
	n.To = types.SetValueMust(types.StringType, nil)
	n.Cc = types.SetValueMust(types.StringType, nil)
	n.Bcc = types.SetValueMust(types.StringType, nil)
	helpers.WriteFields(ctx, n, notification.GetFields(), notificationFields)
}

func (n *Notification) read(ctx context.Context, diags *diag.Diagnostics) *prowlarr.NotificationResource {
	notification := prowlarr.NewNotificationResource()
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnHealthRestored(n.OnHealthRestored.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetIncludeManualGrabs(n.IncludeManualGrabs.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetImplementation(n.Implementation.ValueString())
	notification.SetConfigContract(n.ConfigContract.ValueString())
	diags.Append(n.Tags.ElementsAs(ctx, &notification.Tags, true)...)
	notification.SetFields(helpers.ReadFields(ctx, n, notificationFields))

	return notification
}
