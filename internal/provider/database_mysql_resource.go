// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseMySqlResource{}
var _ resource.ResourceWithImportState = &DatabaseMySqlResource{}

func NewDatabaseMySqlResource() resource.Resource {
	return &DatabaseMySqlResource{}
}

// DatabaseMySqlResource defines the resource implementation.
type DatabaseMySqlResource struct {
	client *client.DatabasusClient
}

func (r *DatabaseMySqlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_mysql"
}

func (r *DatabaseMySqlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "MySQL Database resource to configure connection details to a database within MariDB instance.\n\nPlease note that the DB must be reachable and credentials must be valid, otherwise terraform will fail creating this resource, since Databasus performs a connection check during registration of the new database configuration.",

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
				MarkdownDescription: "The port number of the MySQL instance (e.g. 5432)",
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

func (r *DatabaseMySqlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseMySqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.DatabaseMySqlResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	result, err := r.client.CreateDatabaseMySql(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating MySQL Database Resource",
			"Could not create MySQL Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new MySQL Database Resource resource")

	// Set state to fully populated data
	client.MapResponseToDatabaseMySqlResourceModel(result, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseMySqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.DatabaseMySqlResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetDatabaseMySql(ctx, data.Id.ValueString())
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
	client.MapResponseToDatabaseMySqlResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMySqlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.DatabaseMySqlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	result, err := r.client.UpdateDatabaseMySql(ctx, data.Id.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating MySQL Database Resource",
			"Could not update MySQL Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	client.MapResponseToDatabaseMySqlResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMySqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.DatabaseMySqlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabaseMySql(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting MySQL Database Resource Configuration", "Could not delete MySQL Database Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *DatabaseMySqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
