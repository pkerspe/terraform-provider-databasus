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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseMongoDbResource{}
var _ resource.ResourceWithImportState = &DatabaseMongoDbResource{}

func NewDatabaseMongoDbResource() resource.Resource {
	return &DatabaseMongoDbResource{}
}

// DatabaseMongoDbResource defines the resource implementation.
type DatabaseMongoDbResource struct {
	client *client.DatabasusClient
}

func (r *DatabaseMongoDbResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_mongodb"
}

func (r *DatabaseMongoDbResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "MongoDB Database resource to configure connection details to a database within MariDB instance.\n\nPlease note that the DB must be reachable and credentials must be valid, otherwise terraform will fail creating this resource, since Databasus performs a connection check during registration of the new database configuration.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name (in the Databasus UI) of the database configuration to be created",
				Required:            true,
			},
			"auth_database": schema.StringAttribute{
				MarkdownDescription: "The auth database to use for authentication. E.g. admin",
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
				MarkdownDescription: "The port number of the MongoDB instance (e.g. 27017)",
				Required:            true,
			},
			"is_https": schema.BoolAttribute{
				MarkdownDescription: "Use HTTPS / TLS or not when connecting to the DB host",
				Required:            true,
			},
			"is_direct_connection": schema.BoolAttribute{
				MarkdownDescription: "Connect directly to a single server, skipping replica set discovery. Useful when the server is behind a load balancer, proxy or tunnel. Defaults to false",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"is_srv": schema.BoolAttribute{
				MarkdownDescription: "Enable for MongoDB Atlas SRV connections (mongodb+srv://). Port is not required for SRV connections. Defaults to false",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"cpu_count": schema.Int32Attribute{
				MarkdownDescription: "Number of CPU cores to use for backup and restore operations. Higher values may speed up operations but use more resources. Defaults to 1",
				Optional:            true,
				Default:             int32default.StaticInt32(1),
				Computed:            true,
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

func (r *DatabaseMongoDbResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseMongoDbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.DatabaseMongoDbResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	result, err := r.client.CreateDatabaseMongoDb(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating MongoDB Database Resource",
			"Could not create MongoDB Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new MongoDB Database Resource resource")

	// Set state to fully populated data
	client.MapResponseToDatabaseMongoDbResourceModel(result, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseMongoDbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.DatabaseMongoDbResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetDatabaseMongoDb(ctx, data.Id.ValueString())
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
	client.MapResponseToDatabaseMongoDbResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMongoDbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.DatabaseMongoDbResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	result, err := r.client.UpdateDatabaseMongoDb(ctx, data.Id.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating MongoDB Database Resource",
			"Could not update MongoDB Database Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	client.MapResponseToDatabaseMongoDbResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseMongoDbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.DatabaseMongoDbResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabaseMongoDb(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting MongoDB Database Resource Configuration", "Could not delete MongoDB Database Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *DatabaseMongoDbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
