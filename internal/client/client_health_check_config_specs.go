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
*						Health-Check Config CRUD functions
********************************************************/

type HealthCheckConfigResourceModel struct {
	AttemptsBeforeConsideredAsDown  types.Int32  `tfsdk:"attempts_before_considered_down"`
	DatabaseId                      types.String `tfsdk:"database_id"`
	IntervalMinutes                 types.Int32  `tfsdk:"interval_minutes"`
	HealthCheckEnabled              types.Bool   `tfsdk:"health_check_enabled"`
	SentNotificationWhenUnavailable types.Bool   `tfsdk:"sent_notification_when_unavailable"`
	StoreAttemptsDays               types.Int32  `tfsdk:"store_attempts_days"`
}

func MapResponseToHealthCheckConfigResourceModel(response *HealthCheckConfigResponseModel, data *HealthCheckConfigResourceModel) {
	data.AttemptsBeforeConsideredAsDown = types.Int32Value(response.AttemptsBeforeConsideredAsDown)
	data.DatabaseId = types.StringValue(response.DatabaseId)
	data.IntervalMinutes = types.Int32Value(response.IntervalMinutes)
	data.HealthCheckEnabled = types.BoolValue(response.IsHealthCheckEnabled)
	data.SentNotificationWhenUnavailable = types.BoolValue(response.IsSentNotificationWhenUnavailable)
	data.StoreAttemptsDays = types.Int32Value(response.StoreAttemptsDays)

}

type HealthCheckConfigResponseModel struct {
	AttemptsBeforeConsideredAsDown    int32  `json:"attemptsBeforeConcideredAsDown"`
	DatabaseId                        string `json:"databaseId"`
	IntervalMinutes                   int32  `json:"intervalMinutes"`
	IsHealthCheckEnabled              bool   `json:"isHealthCheckEnabled"`
	IsSentNotificationWhenUnavailable bool   `json:"isSentNotificationWhenUnavailable"`
	StoreAttemptsDays                 int32  `json:"storeAttemptsDays"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallHealthCheckConfigResourceModel(data HealthCheckConfigResourceModel) map[string]any {
	body := map[string]any{
		"attemptsBeforeConcideredAsDown":    data.AttemptsBeforeConsideredAsDown.ValueInt32(),
		"databaseId":                        data.DatabaseId.ValueString(),
		"intervalMinutes":                   data.IntervalMinutes.ValueInt32(),
		"isHealthCheckEnabled":              data.HealthCheckEnabled.ValueBool(),
		"isSentNotificationWhenUnavailable": data.SentNotificationWhenUnavailable.ValueBool(),
		"storeAttemptsDays":                 data.StoreAttemptsDays.ValueInt32(),
	}
	return body
}

func (c *DatabasusClient) CreateHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	body := marshallHealthCheckConfigResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/healthcheck-config", bytes.NewBuffer(b), nil) // the REST API does not return the object for the health-check config, but only a success message
	if err != nil {
		return err
	}

	return nil
}

func (c *DatabasusClient) GetHealthCheckConfig(ctx context.Context, databaseId string) (*HealthCheckConfigResponseModel, *ErrorDetails) {
	var result HealthCheckConfigResponseModel

	err := c.doRequest(ctx, "GET", "/healthcheck-config/"+databaseId, nil, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	body := marshallHealthCheckConfigResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/healthcheck-config", bytes.NewBuffer(b), nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete not supported for health-check config in databasus (disable health-check instead).
func (c *DatabasusClient) DeleteHealthCheckConfig(ctx context.Context, data HealthCheckConfigResourceModel) *ErrorDetails {
	data.HealthCheckEnabled = types.BoolValue(false)
	err := c.UpdateHealthCheckConfig(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
