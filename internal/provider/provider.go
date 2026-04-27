// Copyright (c) KerspeP
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
	"github.com/pkerspe/terraform-provider-databasus/internal/provider/datasources"
	"github.com/pkerspe/terraform-provider-databasus/internal/provider/resources"
)

// Ensure DatabasusProvider satisfies various provider interfaces.
var _ provider.Provider = &DatabasusProvider{}

// DatabasusProvider defines the provider implementation.
type DatabasusProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DatabasusProviderModel describes the provider data model.
type DatabasusProviderModel struct {
	BaseUrl  types.String `tfsdk:"baseurl"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
}

func (p *DatabasusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "databasus"
	resp.Version = p.version
}

func (p *DatabasusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"baseurl": schema.StringAttribute{
				MarkdownDescription: "The REST API base URL from the Databasus instance e.g. https://youserver.local/api/v1",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email (username) of the user to be used for executing requests against the Databasus REST API. Must be a user with admin Role.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the user",
				Required:            true,
			},
		},
	}
}

func (p *DatabasusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config DatabasusProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get token from Databasus for API calls
	token, err := client.GetJWT(config.BaseUrl.ValueString(), config.Email.ValueString(), config.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to authenticate against Databasus REST API",
			err.Error(),
		)
		return
	}

	client := client.NewDatabasusClient(config.BaseUrl.ValueString(), token)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DatabasusProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewWorkspaceResource,
	}
}

func (p *DatabasusProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewWorkspaceDataSource,
		datasources.NewAllWorkspacesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DatabasusProvider{
			version: version,
		}
	}
}
