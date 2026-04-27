package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UsersSettingsDataSource{}

func NewUsersSettingsDataSource() datasource.DataSource {
	return &UsersSettingsDataSource{}
}

// UsersSettingsDataSource defines the data source implementation.
type UsersSettingsDataSource struct {
	client *client.DatabasusClient
}

func (d *UsersSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users_settings"
}

func (d *UsersSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get the global users settings from Databasus",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the settings (seems to be static over time)",
				Computed:            true,
			},
			"allow_external_registrations": schema.BoolAttribute{
				MarkdownDescription: "When enabled, new users can register accounts in Databasus. If disabled, new users can only register via invitation",
				Computed:            true,
			},
			"allow_member_invitations": schema.BoolAttribute{
				MarkdownDescription: "When enabled, existing members can invite new users to join Databasus. If not - only admins can invite users.",
				Computed:            true,
			},
			"member_allowed_to_create_workspaces": schema.BoolAttribute{
				MarkdownDescription: "When enabled, members (non-admin users) can create new workspaces. If not - only admins can create workspaces.",
				Computed:            true,
			},
		},
	}
}

func (d *UsersSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.DatabasusClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.DatabasusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *UsersSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data client.SettingsDataSourceModel

	// Read config FIRST
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetUsersSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Map API → Terraform state
	data.Id = types.StringValue(result.ID)
	data.IsAllowExternalRegistrations = types.BoolValue(result.IsAllowExternalRegistrations)
	data.IsAllowMemberInvitations = types.BoolValue(result.IsAllowMemberInvitations)
	data.IsMemberAllowedToCreateWorkspaces = types.BoolValue(result.IsMemberAllowedToCreateWorkspaces)

	// Logging
	tflog.Trace(ctx, "read settings data source")

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
