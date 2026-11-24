// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NotifierWebhookResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	BodyTemplate  types.String `tfsdk:"body_template"`
	Headers       types.Map `tfsdk:"headers"`
	WebhookMethod types.String `tfsdk:"webhook_method"`
	WebhookUrl    types.String `tfsdk:"webhook_url"`
	LastSaveError types.String `tfsdk:"last_save_error"`
}

func MapResponseToNotifierWebhookResourceModel(ctx context.Context, response *NotifierWebhookResponseModel, data *NotifierWebhookResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.BodyTemplate = types.StringValue(response.Settings.BodyTemplate)
	data.WebhookMethod = types.StringValue(response.Settings.WebhookMethod)
	data.WebhookUrl = types.StringValue(response.Settings.WebhookUrl)
	data.LastSaveError = types.StringValue(response.LastSaveError)
}

type NotifierWebhookResponseModel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	WorkspaceId   string `json:"workspaceId"`
	Settings      NotifierWebhookDetailsResponseModel `json:"settings"`
	LastSaveError string `json:"lastSaveError"`
}

type NotifierWebhookDetailsResponseModel struct {
	Type               string `json:"type"`
	Headers            []NotifierWebhookHeader `json:"headers"`
	BodyTemplate       string `json:"bodyTemplate"`
	WebhookMethod      string `json:"webhookMethod"`
	WebhookUrl         string `json:"webhookUrl"`
	AcceptedCategories []string `json:"acceptedCategories"`
}

type NotifierWebhookHeader struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

func marshallNotifierWebhookResourceModel(ctx context.Context, data NotifierWebhookResourceModel) map[string]any {
	headersMap := make(map[string]string)
	if !data.Headers.IsNull() && !data.Headers.IsUnknown() {
		_ = data.Headers.ElementsAs(ctx, &headersMap, false)
	}
	headerList := make([]map[string]string, 0, len(headersMap))
	for name, value := range headersMap {
		headerList = append(headerList, map[string]string{"name": name, "value": value})
	}
	return map[string]any{
		"name": data.Name.ValueString(),
		"workspaceId": data.WorkspaceId.ValueString(),
		"settings": map[string]any{
			"type": "WEBHOOK",
			"webhookUrl": data.WebhookUrl.ValueString(),
			"webhookMethod": data.WebhookMethod.ValueString(),
			"headers": headerList,
			"bodyTemplate": data.BodyTemplate.ValueString(),
			"acceptedCategories": []string{"ALL"},
		},
	}
}

func (c *DatabasusClient) CreateNotifierWebhook(ctx context.Context, data NotifierWebhookResourceModel) (*NotifierWebhookResponseModel, error) {
	var result NotifierWebhookResponseModel
	b, _ := json.Marshal(marshallNotifierWebhookResourceModel(ctx, data))
	if err := c.doRequest(ctx, "POST", "/notifiers", bytes.NewBuffer(b), &result); err != nil { return nil, err }
	return &result, nil
}

func (c *DatabasusClient) GetNotifierWebhook(ctx context.Context, id string) (*NotifierWebhookResponseModel, *ErrorDetails) {
	var result NotifierWebhookResponseModel
	if err := c.doRequest(ctx, "GET", "/notifiers/"+id, nil, &result); err != nil { return nil, err }
	return &result, nil
}

func (c *DatabasusClient) UpdateNotifierWebhook(ctx context.Context, id string, data NotifierWebhookResourceModel) (*NotifierWebhookResponseModel, error) {
	var result NotifierWebhookResponseModel
	body := marshallNotifierWebhookResourceModel(ctx, data)
	delete(body, "workspaceId")
	b, _ := json.Marshal(body)
	if err := c.doRequest(ctx, "PUT", "/notifiers/"+id, bytes.NewBuffer(b), &result); err != nil { return nil, err }
	return &result, nil
}

func (c *DatabasusClient) DeleteNotifierWebhook(ctx context.Context, id string) *ErrorDetails {
	return c.doRequest(ctx, "DELETE", "/notifiers/"+id, nil, nil)
}
