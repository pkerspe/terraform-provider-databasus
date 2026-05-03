// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseMariaDbResource{}
var _ resource.ResourceWithImportState = &DatabaseMariaDbResource{}

func NewDatabaseMariaDbResource() resource.Resource {
	return &DatabaseMariaDbResource{}
}

// DatabaseMariaDbResource defines the resource implementation.
type DatabaseMariaDbResource struct {
	client *client.DatabasusClient
}

func (r *DatabaseMariaDbResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_mariadb"
}

func (r *DatabaseMariaDbResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "MariaDB Database resource to configure connection details to a database within MariDB instance.\n\nPlease note that the DB must be reachable and credentials must be valid, otherwise terraform will fail creating this resource, since Databasus performs a connection check during registration of the new database configuration.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name (in the Databasus UI) of the database configuration to be created",
				Required:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "The name of the database to be backed up",
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The hostname of the database instance",
				Required:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "The port number of the MariaDB instance (e.g. 5432)",
				Required:            true,
			},
			"is_https": schema.BoolAttribute{
				MarkdownDescription: "Use HTTPS / TLS or not when connecting to the DB host",
				Required:            true,
			},
			"username": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The DB username (role) to be used for creating the backups",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for the DB role to use when creating backups",
			},
			"exclude_events": schema.BoolAttribute{
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				MarkdownDescription: "Skip backing up database events. Enable this if the event scheduler is disabled on your MariaDB server. Defaults to true",
			},
			"workspace_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The id of the workspace the storage belongs to",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The internal Id of the Storage Configuration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DatabaseMariaDbResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseMariaDbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.DatabaseMariaDbResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	result, err := r.client.CreateDatabaseMariaDb(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating MariaDB Database Resource",
			"Could not create MariaDB Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new MariaDB Database Resource resource")

	// Set state to fully populated data
	client.MapResponseToDatabaseMariaDbResourceModel(result, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseMariaDbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.DatabaseMariaDbResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetDatabaseMariaDb(ctx, data.Id.ValueString())
	if err != nil {
		// The Databasus RETS API returns currently an wrong RC 400 in case of te Record not found
		// Terraform expects an empty return in that case without an error in the diagnostics
		// see also https://github.com/databasus/databasus/issues/529
		if err.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Set state to fully populated data
	client.MapResponseToDatabaseMariaDbResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMariaDbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.DatabaseMariaDbResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	result, err := r.client.UpdateDatabaseMariaDb(ctx, data.Id.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating MariaDB Database Resource",
			"Could not update MariaDB Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	client.MapResponseToDatabaseMariaDbResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMariaDbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.DatabaseMariaDbResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabaseMariaDb(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting MariaDB Database Resource Configuration", "Could not delete MariaDB Database Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *DatabaseMariaDbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
