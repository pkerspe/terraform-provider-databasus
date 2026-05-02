// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

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
var _ datasource.DataSource = &WorkspaceDataSource{}

func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

// WorkspaceDataSource defines the data source implementation.
type WorkspaceDataSource struct {
	client *client.DatabasusClient
}

func (d *WorkspaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (d *WorkspaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get information about a configured Workspace in Databasus",

		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when this workspace was created",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workspace",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The workspace Id to fetch details for",
				Required:            true,
			},
		},
	}
}

func (d *WorkspaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data client.WorkspaceDataSourceModel

	// Read config FIRST
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate input
	if data.Id.IsNull() || data.Id.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing Workspace Id",
			"The data source requires an id to query the API.",
		)
		return
	}

	result, err := d.client.GetWorkspace(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Map API → Terraform state
	data.Id = types.StringValue(result.Id)
	data.Name = types.StringValue(result.Name)
	data.CreatedAt = types.StringValue(result.CreatedAt)

	// Logging
	tflog.Trace(ctx, "read workspace data source")

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
