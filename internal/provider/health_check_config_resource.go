// Copyright pkerspe 2026
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HealthCheckConfigResource{}
var _ resource.ResourceWithImportState = &HealthCheckConfigResource{}

func NewHealthCheckConfigResource() resource.Resource {
	return &HealthCheckConfigResource{}
}

// HealthCheckConfigResource defines the resource implementation.
type HealthCheckConfigResource struct {
	client *client.DatabasusClient
}

func (r *HealthCheckConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_health_check_config"
}

func (r *HealthCheckConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource to manage health-check settings for a configured database",

		Attributes: map[string]schema.Attribute{
			"database_id": schema.StringAttribute{
				MarkdownDescription: "The Id of the database config this health-check configuration applies for",
				Required:            true,
			},
			"health_check_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable or disable health-check monitoring for the referenced database",
				Required:            true,
			},
			"sent_notification_when_unavailable": schema.BoolAttribute{
				MarkdownDescription: "Send notifications when database becomes unavailable. Defaults to false.\n\nNote: please make sure you have a notifier configured as well for the database if you set this to true.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"attempts_before_considered_down": schema.Int32Attribute{
				MarkdownDescription: "Number of failed health-check attempts before marking database as down. Defaults to 3",
				Optional:            true,
				Default:             int32default.StaticInt32(3),
				Computed:            true,
			},
			"store_attempts_days": schema.Int32Attribute{
				MarkdownDescription: "How many days to store health check attempt history. Defaults to 7",
				Optional:            true,
				Default:             int32default.StaticInt32(7),
				Computed:            true,
			},
			"interval_minutes": schema.Int32Attribute{
				MarkdownDescription: "How often to check database health (in minutes). Defaults to 1",
				Optional:            true,
				Default:             int32default.StaticInt32(1),
				Computed:            true,
			},
		},
	}
}

func (r *HealthCheckConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HealthCheckConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.HealthCheckConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	err := r.client.CreateHealthCheckConfig(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Health-Check-Config Resource",
			"Could not create Health-Check-Config Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new Health-Check-Config resource")

	// since the databasus REST API does not respond with an updated object we just write the plan to the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HealthCheckConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.HealthCheckConfigResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.DatabaseId.IsNull() || data.DatabaseId.ValueString() == "" || data.DatabaseId.IsUnknown() {
		resp.Diagnostics.AddError(
			"Failed to get health check config",
			"DatabaseId is null or empty. This value is required to fetch the health check configuration.",
		)
		return
	}

	result, err := r.client.GetHealthCheckConfig(ctx, data.DatabaseId.ValueString())
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
	client.MapResponseToHealthCheckConfigResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HealthCheckConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.HealthCheckConfigResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateHealthCheckConfig(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Health-Check-Config Resource",
			"Could not update Health-Check-Config Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// no update is sent from the databasus API in the response, so we just write the same data back to the state (might not even be needed)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HealthCheckConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.HealthCheckConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHealthCheckConfig(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Health-Check-Config Resource Configuration", "Could not delete Health-Check-Config Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *HealthCheckConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
