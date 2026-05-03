// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*******************************************************
*	Backup Config CRUD functions and model definitions
********************************************************/

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
	data.StorageId = types.StringValue(response.Storage.Id)
	data.Enabled = types.BoolValue(response.Enabled)
	data.Interval = types.StringValue(response.BackupInterval.Interval)
	data.TimeOfDay = types.StringValue(response.BackupInterval.TimeOfDay)
	data.Weekday = types.Int32Value(response.BackupInterval.Weekday)
	data.DayOfMonth = types.Int32Value(response.BackupInterval.DayOfMonth)
	data.CronExpression = types.StringValue(response.BackupInterval.CronExpression)
	data.MaxFailedRetryCount = types.Int32Value(response.MaxFailedRetryCount)
	data.Encryption = types.BoolValue(response.Encryption != "NONE")
	data.RetentionPolicyType = types.StringValue(response.RetentionPolicyType)
	data.RetentionCount = types.Int32Value(response.RetentionCount)
	data.RetentionTimePeriod = types.StringValue(response.RetentionTimePeriod)
	data.RetentionGfsHours = types.Int32Value(response.RetentionGfsHours)
	data.RetentionGfsDays = types.Int32Value(response.RetentionGfsDays)
	data.RetentionGfsWeeks = types.Int32Value(response.RetentionGfsWeeks)
	data.RetentionGfsMonths = types.Int32Value(response.RetentionGfsMonths)
	data.RetentionGfsYears = types.Int32Value(response.RetentionGfsYears)
	data.SendNotificationsOnBackupSuccess = types.BoolValue(containsString(response.SendNotificationsOn, "BACKUP_SUCCESS"))
	data.SendNotificationsOnBackupFailure = types.BoolValue(containsString(response.SendNotificationsOn, "BACKUP_FAILED"))
}

type BackupConfigResponseModel struct {
	BackupInterval      BackupIntervalDetailsResponseModel `json:"backupInterval"`
	DatabaseId          string                             `json:"databaseId"`
	Enabled             bool                               `json:"isBackupsEnabled"`
	MaxFailedRetryCount int32                              `json:"maxFailedTriesCount"`
	IsRetryIfFailed     bool                               `json:"isRetryIfFailed"`
	Encryption          string                             `json:"encryption"`
	RetentionPolicyType string                             `json:"retentionPolicyType"`
	RetentionCount      int32                              `json:"retentionCount"`
	RetentionTimePeriod string                             `json:"retentionTimePeriod"`
	RetentionGfsHours   int32                              `json:"retentionGfsHours"`
	RetentionGfsDays    int32                              `json:"retentionGfsDays"`
	RetentionGfsWeeks   int32                              `json:"retentionGfsWeeks"`
	RetentionGfsMonths  int32                              `json:"retentionGfsMonths"`
	RetentionGfsYears   int32                              `json:"retentionGfsYears"`
	SendNotificationsOn []string                           `json:"sendNotificationsOn"`
	Storage             BackupStorageDetailsResponseModel  `json:"storage"`
}

type BackupIntervalDetailsResponseModel struct {
	Interval       string `json:"interval"`
	TimeOfDay      string `json:"timeOfDay"`
	Weekday        int32  `json:"weekday"`
	DayOfMonth     int32  `json:"dayOfMonth"`
	CronExpression string `json:"cronExpression"`
}

type BackupStorageDetailsResponseModel struct {
	Id string `json:"id"`
}

func containsString(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func ternary(cond bool, a, b any) any {
	if cond {
		return a
	}
	return b
}

// internal helper to transform Model to map that can be used in request body.
func marshallBackupConfigResourceModel(data BackupConfigResourceModel) map[string]any {
	var notificationOptions []string

	if data.SendNotificationsOnBackupSuccess.ValueBool() {
		notificationOptions = append(notificationOptions, "BACKUP_SUCCESS")
	}
	if data.SendNotificationsOnBackupFailure.ValueBool() {
		notificationOptions = append(notificationOptions, "BACKUP_FAILED")
	}

	body := map[string]any{
		"backupInterval": map[string]any{
			"cronExpression": data.CronExpression.ValueString(),
			"dayOfMonth":     data.DayOfMonth.ValueInt32(),
			"interval":       data.Interval.ValueString(),
			"weekday":        data.Weekday.ValueInt32(),
			"timeOfDay":      data.TimeOfDay.ValueString(),
		},
		"databaseId":          data.DatabaseId.ValueString(),
		"encryption":          ternary(data.Encryption.ValueBool(), "ENCRYPTED", "NONE"),
		"isBackupsEnabled":    data.Enabled.ValueBool(),
		"isRetryIfFailed":     ternary(data.MaxFailedRetryCount.ValueInt32() > 0, true, false),
		"maxFailedTriesCount": data.MaxFailedRetryCount.ValueInt32(),
		"retentionCount":      data.RetentionCount.ValueInt32(),
		"retentionGfsDays":    data.RetentionGfsDays.ValueInt32(),
		"retentionGfsHours":   data.RetentionGfsHours.ValueInt32(),
		"retentionGfsMonths":  data.RetentionGfsMonths.ValueInt32(),
		"retentionGfsWeeks":   data.RetentionGfsWeeks.ValueInt32(),
		"retentionGfsYears":   data.RetentionGfsYears.ValueInt32(),
		"retentionPolicyType": data.RetentionPolicyType.ValueString(),
		"retentionTimePeriod": data.RetentionTimePeriod.ValueString(),
		"sendNotificationsOn": notificationOptions,
		"storage": map[string]any{
			"id": data.StorageId.ValueString(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateBackupConfig(ctx context.Context, data BackupConfigResourceModel) (*BackupConfigResponseModel, error) {
	var result BackupConfigResponseModel
	body := marshallBackupConfigResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/backup-configs/save", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetBackupConfig(ctx context.Context, databaseId string) (resultModel *BackupConfigResponseModel, errorD *ErrorDetails) {
	var result BackupConfigResponseModel

	err := c.doRequest(ctx, "GET", "/backup-configs/database/"+databaseId, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateBackupConfig(ctx context.Context, data BackupConfigResourceModel) (*BackupConfigResponseModel, error) {
	// API provides only a single endpoint for create and update
	return c.CreateBackupConfig(ctx, data)
}

// Delete is not possible for backup configs, can only set to false, but need to use update then.
func (c *DatabasusClient) DeleteBackupConfig(ctx context.Context) error {
	return nil
}
