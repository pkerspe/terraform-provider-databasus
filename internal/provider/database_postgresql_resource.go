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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabasePostgresqlResource{}
var _ resource.ResourceWithImportState = &DatabasePostgresqlResource{}

func NewDatabasePostgresqlResource() resource.Resource {
	return &DatabasePostgresqlResource{}
}

// DatabasePostgresqlResource defines the resource implementation.
type DatabasePostgresqlResource struct {
	client *client.DatabasusClient
}

func (r *DatabasePostgresqlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_postgresql"
}

func (r *DatabasePostgresqlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "PostgreSQL Database resource for Remote Backup",

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
				MarkdownDescription: "The port number of the PostgreSQL instance (e.g. 5432)",
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
			"include_schemas": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "List of Schema names to include in the Backup. Use empty list for all.",
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

func (r *DatabasePostgresqlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabasePostgresqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.DatabasePostgresqlResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	result, err := r.client.CreateDatabasePostgresql(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Postgresql Database Resource",
			"Could not create Postgresql Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new Postgresql Database Resource resource")

	// Set state to fully populated data
	client.MapResponseToDatabasePostgresqlResourceModel(result, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabasePostgresqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.DatabasePostgresqlResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetDatabasePostgresql(ctx, data.Id.ValueString())
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
	client.MapResponseToDatabasePostgresqlResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabasePostgresqlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.DatabasePostgresqlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	result, err := r.client.UpdateDatabasePostgresql(ctx, data.Id.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Postgresql Database Resource",
			"Could not update Postgresql Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	client.MapResponseToDatabasePostgresqlResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabasePostgresqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.DatabasePostgresqlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabasePostgresql(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Postgresql Database Resource Configuration", "Could not delete Postgresql Database Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *DatabasePostgresqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
