// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*******************************************************
* Backup Config CRUD functions and model definitions
*******************************************************/

type BackupConfigResourceModel struct {
	Enabled                          types.Bool   `tfsdk:"enabled"`
	DatabaseId                       types.String `tfsdk:"database_id"`
	StorageId                        types.String `tfsdk:"storage_id"`
	Interval                         types.String `tfsdk:"interval"`
	TimeOfDay                        types.String `tfsdk:"time_of_day"`
	Weekday                          types.Int32  `tfsdk:"weekday"`
	DayOfMonth                       types.Int32  `tfsdk:"day_of_month"`
	CronExpression                   types.String `tfsdk:"cron_expression"`
	MaxFailedRetryCount              types.Int32  `tfsdk:"max_failed_retry_count"`
	Encryption                       types.Bool   `tfsdk:"encryption"`
	RetentionPolicyType              types.String `tfsdk:"retention_policy_type"`
	RetentionCount                   types.Int32  `tfsdk:"retention_count"`
	RetentionTimePeriod              types.String `tfsdk:"retention_time_period"`
	RetentionGfsHours                types.Int32  `tfsdk:"retention_gfs_hours"`
	RetentionGfsDays                 types.Int32  `tfsdk:"retention_gfs_days"`
	RetentionGfsWeeks                types.Int32  `tfsdk:"retention_gfs_weeks"`
	RetentionGfsMonths               types.Int32  `tfsdk:"retention_gfs_months"`
	RetentionGfsYears                types.Int32  `tfsdk:"retention_gfs_years"`
	SendNotificationsOnBackupSuccess types.Bool   `tfsdk:"send_notifications_on_backup_success"`
	SendNotificationsOnBackupFailure types.Bool   `tfsdk:"send_notifications_on_backup_failure"`
}

func MapResponseToBackupConfigResourceModel(response *BackupConfigResponseModel, data *BackupConfigResourceModel) {
	data.DatabaseId = types.StringValue(response.DatabaseId)
	if response.StorageId == nil {
		data.StorageId = types.StringValue("")
	} else {
		data.StorageId = types.StringValue(*response.StorageId)
	}
	data.Enabled = types.BoolValue(response.IsBackupsEnabled)
	data.Interval = types.StringValue(response.Settings.Schedule.Type)
	data.TimeOfDay = types.StringValue(response.Settings.Schedule.TimeOfDay)
	data.Weekday = types.Int32Value(response.Settings.Schedule.Weekday)
	data.DayOfMonth = types.Int32Value(response.Settings.Schedule.DayOfMonth)
	data.CronExpression = types.StringValue(response.Settings.Schedule.CronExpression)
	data.MaxFailedRetryCount = types.Int32Value(response.Settings.MaxFailedTriesCount)
	data.Encryption = types.BoolValue(response.Encryption != "NONE")
	data.RetentionPolicyType = types.StringValue(response.Settings.Retention.Policy)
	data.RetentionCount = types.Int32Value(response.Settings.Retention.Count)
	data.RetentionTimePeriod = types.StringValue(response.Settings.Retention.Period)
	data.RetentionGfsHours = types.Int32Value(response.Settings.Retention.GfsHours)
	data.RetentionGfsDays = types.Int32Value(response.Settings.Retention.GfsDays)
	data.RetentionGfsWeeks = types.Int32Value(response.Settings.Retention.GfsWeeks)
	data.RetentionGfsMonths = types.Int32Value(response.Settings.Retention.GfsMonths)
	data.RetentionGfsYears = types.Int32Value(response.Settings.Retention.GfsYears)
	data.SendNotificationsOnBackupSuccess = types.BoolValue(containsString(response.Settings.NotificationTriggers, "BACKUP_SUCCESS"))
	data.SendNotificationsOnBackupFailure = types.BoolValue(containsString(response.Settings.NotificationTriggers, "BACKUP_FAILED"))
}

type BackupConfigResponseModel struct {
	DatabaseId       string                        `json:"databaseId"`
	IsBackupsEnabled bool                          `json:"isBackupsEnabled"`
	StorageId        *string                       `json:"storageId"`
	Encryption       string                        `json:"encryption"`
	Settings         LogicalBackupSettingsResponse `json:"settings"`
}

type LogicalBackupSettingsResponse struct {
	Type                  string                    `json:"type"`
	Schedule              BackupScheduleResponse    `json:"schedule"`
	Retention             LogicalRetentionResponse  `json:"retention"`
	IsRetryIfFailed       bool                      `json:"isRetryIfFailed"`
	MaxFailedTriesCount   int32                     `json:"maxFailedTriesCount"`
	NotificationTriggers  []string                  `json:"notificationTriggers"`
}

type BackupScheduleResponse struct {
	Type           string `json:"type"`
	TimeOfDay      string `json:"timeOfDay"`
	Weekday        int32  `json:"weekday"`
	DayOfMonth     int32  `json:"dayOfMonth"`
	CronExpression string `json:"cronExpression"`
}

type LogicalRetentionResponse struct {
	Policy    string `json:"policy"`
	Period    string `json:"period"`
	Count     int32  `json:"count"`
	GfsHours  int32  `json:"gfsHours"`
	GfsDays   int32  `json:"gfsDays"`
	GfsWeeks  int32  `json:"gfsWeeks"`
	GfsMonths int32  `json:"gfsMonths"`
	GfsYears  int32  `json:"gfsYears"`
}

func containsString(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

func marshallBackupConfigResourceModel(data BackupConfigResourceModel) map[string]any {
	notificationTriggers := make([]string, 0, 2)
	if data.SendNotificationsOnBackupSuccess.ValueBool() {
		notificationTriggers = append(notificationTriggers, "BACKUP_SUCCESS")
	}
	if data.SendNotificationsOnBackupFailure.ValueBool() {
		notificationTriggers = append(notificationTriggers, "BACKUP_FAILED")
	}

	return map[string]any{
		"isBackupsEnabled": data.Enabled.ValueBool(),
		"storageId":        data.StorageId.ValueString(),
		"encryption":       encryptionValue(data.Encryption.ValueBool()),
		"settings": map[string]any{
			"type": "LOGICAL",
			"schedule": map[string]any{
				"type":           data.Interval.ValueString(),
				"timeOfDay":      data.TimeOfDay.ValueString(),
				"weekday":        data.Weekday.ValueInt32(),
				"dayOfMonth":     data.DayOfMonth.ValueInt32(),
				"cronExpression": data.CronExpression.ValueString(),
			},
			"retention": map[string]any{
				"policy":    data.RetentionPolicyType.ValueString(),
				"period":    data.RetentionTimePeriod.ValueString(),
				"count":     data.RetentionCount.ValueInt32(),
				"gfsHours":  data.RetentionGfsHours.ValueInt32(),
				"gfsDays":   data.RetentionGfsDays.ValueInt32(),
				"gfsWeeks":  data.RetentionGfsWeeks.ValueInt32(),
				"gfsMonths": data.RetentionGfsMonths.ValueInt32(),
				"gfsYears":  data.RetentionGfsYears.ValueInt32(),
			},
			"isRetryIfFailed":      data.MaxFailedRetryCount.ValueInt32() > 0,
			"maxFailedTriesCount":  data.MaxFailedRetryCount.ValueInt32(),
			"notificationTriggers": notificationTriggers,
		},
	}
}

func encryptionValue(enabled bool) string {
	if enabled {
		return "ENCRYPTED"
	}
	return "NONE"
}

func (c *DatabasusClient) CreateBackupConfig(ctx context.Context, data BackupConfigResourceModel) (*BackupConfigResponseModel, error) {
	return c.saveBackupConfig(ctx, data)
}

func (c *DatabasusClient) saveBackupConfig(ctx context.Context, data BackupConfigResourceModel) (*BackupConfigResponseModel, error) {
	var result BackupConfigResponseModel
	body, err := json.Marshal(marshallBackupConfigResourceModel(data))
	if err != nil {
		return nil, fmt.Errorf("marshal backup config: %w", err)
	}

	path := "/backup-configs/" + data.DatabaseId.ValueString()
	if err := c.doRequest(ctx, "PUT", path, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetBackupConfig(ctx context.Context, databaseID string) (*BackupConfigResponseModel, *ErrorDetails) {
	var result BackupConfigResponseModel
	path := "/backup-configs/" + databaseID
	if err := c.doRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateBackupConfig(ctx context.Context, data BackupConfigResourceModel) (*BackupConfigResponseModel, error) {
	return c.saveBackupConfig(ctx, data)
}

// Backup configs have no delete endpoint; removing the Terraform resource does not delete the server config.
func (c *DatabasusClient) DeleteBackupConfig(ctx context.Context) error {
	return nil
}
