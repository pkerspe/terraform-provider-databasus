package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UsersSettingsResource{}
var _ resource.ResourceWithImportState = &UsersSettingsResource{}

func NewUsersSettingsResource() resource.Resource {
	return &UsersSettingsResource{}
}

// UsersSettingsResource defines the resource implementation.
type UsersSettingsResource struct {
	client *client.DatabasusClient
}

func (r *UsersSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users_settings"
}

func (r *UsersSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource to manage the global users settings from Databasus",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the settings (seems to be static over time)",
				Computed:            true,
			},
			"allow_external_registrations": schema.BoolAttribute{
				MarkdownDescription: "When enabled, new users can register accounts in Databasus. If disabled, new users can only register via invitation",
				Required:            true,
			},
			"allow_member_invitations": schema.BoolAttribute{
				MarkdownDescription: "When enabled, existing members can invite new users to join Databasus. If not - only admins can invite users.",
				Required:            true,
			},
			"member_allowed_to_create_workspaces": schema.BoolAttribute{
				MarkdownDescription: "When enabled, members (non-admin users) can create new workspaces. If not - only admins can create workspaces.",
				Required:            true,
			},
		},
	}
}

func (r *UsersSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.DatabasusClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.DatabasusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *UsersSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.SettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create (actually just Update) the Global Users Settings
	usersSettings, err := r.client.CreateUsersSettings(ctx,
		plan.IsAllowExternalRegistrations.ValueBool(),
		plan.IsAllowMemberInvitations.ValueBool(),
		plan.IsMemberAllowedToCreateWorkspaces.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating users settings",
			"Could not create users settings, unexpected error: "+err.Error(),
		)
		return
	}

	// save response values into the Terraform state.
	plan.Id = types.StringValue(usersSettings.ID)
	plan.IsAllowExternalRegistrations = types.BoolValue(usersSettings.IsAllowExternalRegistrations)
	plan.IsAllowMemberInvitations = types.BoolValue(usersSettings.IsAllowMemberInvitations)
	plan.IsMemberAllowedToCreateWorkspaces = types.BoolValue(usersSettings.IsMemberAllowedToCreateWorkspaces)

	tflog.Trace(ctx, "created new users settings resource")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UsersSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.SettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetUsersSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Map API → Terraform state
	data.Id = types.StringValue(result.ID)
	data.IsAllowExternalRegistrations = types.BoolValue(result.IsAllowExternalRegistrations)
	data.IsAllowMemberInvitations = types.BoolValue(result.IsAllowMemberInvitations)
	data.IsMemberAllowedToCreateWorkspaces = types.BoolValue(result.IsMemberAllowedToCreateWorkspaces)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsersSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.SettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing users settings
	usersSettings, err := r.client.UpdateUsersSettings(ctx,
		data.Id.ValueString(),
		data.IsAllowExternalRegistrations.ValueBool(),
		data.IsAllowMemberInvitations.ValueBool(),
		data.IsMemberAllowedToCreateWorkspaces.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating global users settings",
			"Could not update users settings, unexpected error: "+err.Error(),
		)
		return
	}

	// Map API response to Terraform state
	data.Id = types.StringValue(usersSettings.ID)
	data.IsAllowExternalRegistrations = types.BoolValue(usersSettings.IsAllowExternalRegistrations)
	data.IsAllowMemberInvitations = types.BoolValue(usersSettings.IsAllowMemberInvitations)
	data.IsMemberAllowedToCreateWorkspaces = types.BoolValue(usersSettings.IsMemberAllowedToCreateWorkspaces)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// note: deleting this resource is not really possible, the client will not perform any action
func (r *UsersSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.SettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUsersSettings(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting gLobal users settings", "Could not delete users settings, unexpected error: "+err.Error())
		return
	}
}

func (r *UsersSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
