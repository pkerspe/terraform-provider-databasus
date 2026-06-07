// Copyright pkerspe 2026
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
var _ datasource.DataSource = &WorkspaceByNameDataSource{}

func NewWorkspaceByNameDataSource() datasource.DataSource {
	return &WorkspaceByNameDataSource{}
}

// WorkspaceByNameDataSource defines the data source implementation.
type WorkspaceByNameDataSource struct {
	client *client.DatabasusClient
}

func (d *WorkspaceByNameDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_by_name"
}

func (d *WorkspaceByNameDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get information about a configured Workspace in Databasus. The name (instead of the ID) is used to query the worksapce in this Data Source. NOTE: this only fetches the first matching workspace for the given name.",

		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when this workspace was created",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the workspace",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The workspace name to fetch details for",
				Required:            true,
			},
		},
	}
}

func (d *WorkspaceByNameDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceByNameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data client.WorkspaceDataSourceModel

	// Read config FIRST
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate input
	if data.Name.IsNull() || data.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing Workspace Name",
			"The data source requires a name to query the API.",
		)
		return
	}

	result, err := d.client.GetWorkspaceByName(ctx, data.Name.ValueString())
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
