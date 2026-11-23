// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HealthCheckConfigResourceModel struct {
	AttemptsBeforeConsideredAsDown types.Int32 `tfsdk:"attempts_before_considered_down"`
	DatabaseId types.String `tfsdk:"database_id"`
	IntervalMinutes types.Int32 `tfsdk:"interval_minutes"`
	HealthCheckEnabled types.Bool `tfsdk:"health_check_enabled"`
	SentNotificationWhenUnavailable types.Bool `tfsdk:"sent_notification_when_unavailable"`
	StoreAttemptsDays types.Int32 `tfsdk:"store_attempts_days"`
}

type HealthCheckConfigResponseModel struct {
	DatabaseId string `json:"databaseId"`
	FailedAttemptsBeforeUnavailable int32 `json:"failedAttemptsBeforeUnavailable"`
	IntervalMinutes int32 `json:"intervalMinutes"`
	IsHealthcheckEnabled bool `json:"isHealthcheckEnabled"`
	IsNotificationEnabled bool `json:"isNotificationEnabled"`
	AttemptRetentionDays int32 `json:"attemptRetentionDays"`
	ShouldConfirmEnabling bool `json:"shouldConfirmEnabling"`
}

func MapResponseToHealthCheckConfigResourceModel(response *HealthCheckConfigResponseModel, data *HealthCheckConfigResourceModel) {
	data.AttemptsBeforeConsideredAsDown = types.Int32Value(response.FailedAttemptsBeforeUnavailable)
	data.DatabaseId = types.StringValue(response.DatabaseId)
	data.IntervalMinutes = types.Int32Value(response.IntervalMinutes)
	data.HealthCheckEnabled = types.BoolValue(response.IsHealthcheckEnabled)
	data.SentNotificationWhenUnavailable = types.BoolValue(response.IsNotificationEnabled)
	data.StoreAttemptsDays = types.Int32Value(response.AttemptRetentionDays)
}

func marshallHealthCheckConfigResourceModel(data HealthCheckConfigResourceModel) map[string]any {
	return map[string]any{
		"isHealthcheckEnabled": data.HealthCheckEnabled.ValueBool(),
		"isNotificationEnabled": data.SentNotificationWhenUnavailable.ValueBool(),
		"intervalMinutes": data.IntervalMinutes.ValueInt32(),
		"failedAttemptsBeforeUnavailable": data.AttemptsBeforeConsideredAsDown.ValueInt32(),
		"attemptRetentionDays": data.StoreAttemptsDays.ValueInt32(),
	}
}

func (c *DatabasusClient) CreateHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	return c.UpdateHealthCheckConfig(ctx, data)
}

func (c *DatabasusClient) GetHealthCheckConfig(ctx context.Context, databaseId string) (*HealthCheckConfigResponseModel, *ErrorDetails) {
	var result HealthCheckConfigResponseModel
	if err := c.doRequest(ctx, "GET", "/healthcheck-configs/"+databaseId, nil, &result); err != nil { return nil, err }
	return &result, nil
}

func (c *DatabasusClient) UpdateHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	body, _ := json.Marshal(marshallHealthCheckConfigResourceModel(data))
	var result HealthCheckConfigResponseModel
	if err := c.doRequest(ctx, "PUT", "/healthcheck-configs/"+data.DatabaseId.ValueString(), bytes.NewBuffer(body), &result); err != nil { return err }
	return nil
}

func (c *DatabasusClient) DeleteHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	data.HealthCheckEnabled = types.BoolValue(false)
	return c.UpdateHealthCheckConfig(ctx, data)
}
