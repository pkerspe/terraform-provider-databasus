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
*	Webhook Notifier CRUD functions and model definitions
********************************************************/

type NotifierWebhookResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	BodyTemplate  types.String `tfsdk:"body_template"`
	Headers       types.Map    `tfsdk:"headers"`
	WebhookMethod types.String `tfsdk:"webhook_method"`
	WebhookUrl    types.String `tfsdk:"webhook_url"`
	LastSaveError types.String `tfsdk:"last_save_error"`
}

func MapResponseToNotifierWebhookResourceModel(ctx context.Context, response *NotifierWebhookResponseModel, data *NotifierWebhookResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.BodyTemplate = types.StringValue(response.WebhookNotifier.BodyTemplate)
	data.WebhookMethod = types.StringValue(response.WebhookNotifier.WebhookMethod)
	data.WebhookUrl = types.StringValue(response.WebhookNotifier.WebhookUrl)
	data.LastSaveError = types.StringValue(response.LastSaveError)

	// convert header array to a map structure
	// headersMap := headersSliceToMap(response.WebhookNotifier.Headers)
	// tfMap, _ := types.MapValueFrom(ctx, types.StringType, headersMap)
	// data.Headers = tfMap
}

type NotifierWebhookResponseModel struct {
	Id              string                              `json:"id"`
	Name            string                              `json:"name"`
	WorkspaceId     string                              `json:"workspaceId"`
	WebhookNotifier NotifierWebhookDetailsResponseModel `json:"webhookNotifier"`
	LastSaveError   string                              `json:"lastSaveError"`
}

type NotifierWebhookDetailsResponseModel struct {
	NotifierId    string   `json:"notifierId"`
	Headers       []Header `json:"headers"`
	BodyTemplate  string   `json:"bodyTemplate"`
	WebhookMethod string   `json:"webhookMethod"`
	WebhookUrl    string   `json:"webhookUrl"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallNotifierWebhookResourceModel(ctx context.Context, data NotifierWebhookResourceModel) map[string]any {
	headers := make(map[string]string)
	data.Headers.ElementsAs(ctx, &headers, false)
	headerList := make([]Header, 0, len(headers))
	for k, v := range headers {
		headerList = append(headerList, Header{
			Key:   k,
			Value: v,
		})
	}

	body := map[string]any{
		"notifierType": "WEBHOOK",
		"name":         data.Name.ValueString(),
		"workspaceId":  data.WorkspaceId.ValueString(),
		"webhookNotifier": map[string]any{
			"webhookUrl":    data.WebhookUrl.ValueString(),
			"webhookMethod": data.WebhookMethod.ValueString(),
			"headers":       headerList,
			"bodyTemplate":  data.BodyTemplate.ValueString(),
			//"notifierId"
		},
	}
	return body
}

func (c *DatabasusClient) CreateNotifierWebhook(ctx context.Context, data NotifierWebhookResourceModel) (*NotifierWebhookResponseModel, error) {
	var result NotifierWebhookResponseModel
	body := marshallNotifierWebhookResourceModel(ctx, data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/notifiers", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetNotifierWebhook(ctx context.Context, id string) (resultModel *NotifierWebhookResponseModel, errorD *ErrorDetails) {
	var result NotifierWebhookResponseModel

	err := c.doRequest(ctx, "GET", "/notifiers/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateNotifierWebhook(ctx context.Context, id string, data NotifierWebhookResourceModel) (*NotifierWebhookResponseModel, error) {
	var result NotifierWebhookResponseModel

	body := marshallNotifierWebhookResourceModel(ctx, data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/notifiers", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the configuration.
func (c *DatabasusClient) DeleteNotifierWebhook(ctx context.Context, id string) *ErrorDetails {
	errorD := c.doRequest(ctx, "DELETE", "/notifiers/"+id, nil, nil)
	if errorD != nil {
		return errorD
	}
	return nil
}
