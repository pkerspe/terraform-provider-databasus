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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkerspe/terraform-provider-databasus/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BackupConfigResource{}
var _ resource.ResourceWithImportState = &BackupConfigResource{}

func NewBackupConfigResource() resource.Resource {
	return &BackupConfigResource{}
}

// BackupConfigResource defines the resource implementation.
type BackupConfigResource struct {
	client *client.DatabasusClient
}

func (r *BackupConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_config"
}

func (r *BackupConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Backup Configuration resource for define the backup job details like intervals and the target storage",

		Attributes: map[string]schema.Attribute{
			"database_id": schema.StringAttribute{
				MarkdownDescription: "The id of the database config to create the backup config for",
				Required:            true,
			},
			"interval": schema.StringAttribute{
				MarkdownDescription: "The schedule interval for the backup. Can be HOURLY, DAILY, WEEKLY, MONTHLY",
				Required:            true,
			},
			"time_of_day": schema.StringAttribute{
				MarkdownDescription: "Only needed when interval is set to DAILY, WEEKLY, MONTHLY or CRON. The time of the day in 24 hour format when to perform the backup. Defaults to 00:00 UTC",
				Default:             stringdefault.StaticString("00:00"),
				Computed:            true,
				Optional:            true,
			},
			"weekday": schema.Int32Attribute{
				MarkdownDescription: "Only needed when interval is set to WEEKLY. The weekday when to execute the backup, Defaults to 1 for Monday",
				Default:             int32default.StaticInt32(1),
				Computed:            true,
				Optional:            true,
			},
			"day_of_month": schema.Int32Attribute{
				MarkdownDescription: "Only needed when interval is set to MONTHLY. The day of the month when to execute the backup. Defaults to 1",
				Default:             int32default.StaticInt32(1),
				Computed:            true,
				Optional:            true,
			},
			"cron_expression": schema.StringAttribute{
				MarkdownDescription: "Only needed when interval is set to CRON. Cron format: minute hour day month weekday (UTC). Defaults to 0 0 * * * (Daily at Midnight UTC)",
				Default:             stringdefault.StaticString("0 0 * * *"),
				Computed:            true,
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Flag to enable and disable backups completely. Defaults to true",
				Default:             booldefault.StaticBool(true),
				Computed:            true,
				Optional:            true,
			},
			"max_failed_retry_count": schema.Int32Attribute{
				MarkdownDescription: "Flag to configure if and how many failed backups should be retried. Set to 0 to disable retries. Default to 0",
				Default:             int32default.StaticInt32(0),
				Computed:            true,
				Optional:            true,
			},
			"encryption": schema.BoolAttribute{
				MarkdownDescription: "Flag to configure encryption of backups. If backup is encrypted, backup files in your storage (S3, local, etc.) cannot be used directly. You can restore backups through Databasus or download them unencrypted via the 'Download' button in the UI. Defaults to true",
				Default:             booldefault.StaticBool(true),
				Computed:            true,
				Optional:            true,
			},
			"retention_policy_type": schema.StringAttribute{
				MarkdownDescription: "Retention policy type to be used. Valid values are: TIME_PERIOD (time based backup policy, keeping backup of last periods defined in retention_time_period), COUNT (Keep only the specified number of most recent backups. Older backups beyond this count are automatically deleted, count is defined with retention_count), GFS (Grandfather-Father-Son rotation: keep the last N hourly, daily, weekly, monthly and yearly backups. This allows keeping backups over long periods of time within a reasonable storage space.)",
				Required:            true,
			},
			"retention_count": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to COUNT. Number of last backups to keep. Defaults to 14",
				Default:             int32default.StaticInt32(14),
				Computed:            true,
				Optional:            true,
			},
			"retention_time_period": schema.StringAttribute{
				MarkdownDescription: "Only used if Retention Policy is set to TIME_PERIOD. How long to keep the backups. Backups older than this period are automatically deleted. Allowed values: DAY, WEEK, MONTH, 3_MONTH, 6_MONTH, YEAR, 2_YEARS, 3_YEARS. Defaults to MONTH",
				Default:             stringdefault.StaticString("MONTH"),
				Computed:            true,
				Optional:            true,
			},
			"retention_gfs_hours": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to GFS. Number of Daily backups. Default to 24",
				Default:             int32default.StaticInt32(24),
				Computed:            true,
				Optional:            true,
			},
			"retention_gfs_days": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to GFS. Number of Daily backups. Defaults to 14",
				Default:             int32default.StaticInt32(14),
				Computed:            true,
				Optional:            true,
			},
			"retention_gfs_weeks": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to GFS. Number of Weekly backups. Defaults to 4",
				Default:             int32default.StaticInt32(4),
				Computed:            true,
				Optional:            true,
			},
			"retention_gfs_months": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to GFS. Number of Monthly backups. Defaults to 12",
				Default:             int32default.StaticInt32(12),
				Computed:            true,
				Optional:            true,
			},
			"retention_gfs_years": schema.Int32Attribute{
				MarkdownDescription: "Only used if Retention Policy is set to GFS. Number of Yearly backups. Defaults to 3",
				Default:             int32default.StaticInt32(12),
				Computed:            true,
				Optional:            true,
			},
			"send_notifications_on_backup_success": schema.BoolAttribute{
				MarkdownDescription: "Flag to indicate if notifications should be send when a backup completed successful. Defaults to false",
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Optional:            true,
			},
			"send_notifications_on_backup_failure": schema.BoolAttribute{
				MarkdownDescription: "Flag to indicate if notifications should be send when a backup fails. Defaults to true",
				Default:             booldefault.StaticBool(true),
				Computed:            true,
				Optional:            true,
			},
			"storage_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The id of storage to use for the backups",
			},
		},
	}
}

func (r *BackupConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan client.BackupConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Configuration
	result, err := r.client.CreateBackupConfig(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Backup Configuration Resource",
			"Could not create Backup Configuration Resource, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created new Backup Configuration Resource resource")

	// Set state to fully populated data
	client.MapResponseToBackupConfigResourceModel(result, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BackupConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data client.BackupConfigResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetBackupConfig(ctx, data.DatabaseId.ValueString())
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
	client.MapResponseToBackupConfigResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data client.BackupConfigResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	result, err := r.client.UpdateBackupConfig(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Backup Configuration Resource",
			"Could not update Backup Configuration Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	client.MapResponseToBackupConfigResourceModel(result, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state client.BackupConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBackupConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Backup Configuration Resource Configuration", "Could not delete Backup Configuration Resource, unexpected error: "+err.Error())
		return
	}
}

func (r *BackupConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
