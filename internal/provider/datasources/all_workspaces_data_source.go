// Copyright (c) KerspeP
// SPDX-License-Identifier: Apache-2.0

package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

var _ datasource.DataSource = &AllWorkspacesDataSource{}

func NewAllWorkspacesDataSource() datasource.DataSource {
	return &AllWorkspacesDataSource{}
}

type AllWorkspacesDataSource struct {
	client *client.DatabasusClient
}

type AllWorkspacesDataSourceModel struct {
	Workspaces []client.WorkspaceDataSourceModel `tfsdk:"workspaces"`
}

func (d *AllWorkspacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_workspaces"
}

func (d *AllWorkspacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get a list of all configured Workspaces in Databasus",

		Attributes: map[string]schema.Attribute{
			"workspaces": schema.ListNestedAttribute{
				Description: "List of all existing workspaces.",
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the workspace.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the workspace.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the workspace was created.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AllWorkspacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AllWorkspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllWorkspacesDataSourceModel

	// Read config FIRST
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API
	results, err := d.client.ListWorkspaces(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Map API → Terraform
	workspaces := make([]client.WorkspaceDataSourceModel, 0, len(results.Items))

	for _, w := range results.Items {
		workspaces = append(workspaces, client.WorkspaceDataSourceModel{
			Id:        types.StringValue(w.Id),
			Name:      types.StringValue(w.Name),
			CreatedAt: types.StringValue(w.CreatedAt),
		})
	}

	data.Workspaces = workspaces

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
