// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageLocalResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	LastSaveError types.String `tfsdk:"last_save_error"`
}

func MapResponseToStorageLocalResourceModel(response *StorageLocalResponseModel, data *StorageLocalResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.LastSaveError = types.StringValue(response.LastSaveError)
}

type StorageLocalResponseModel struct {
	Id            string                           `json:"id"`
	Name          string                           `json:"name"`
	WorkspaceId   string                           `json:"workspaceId"`
	Settings      StorageLocalDetailsResponseModel `json:"settings"`
	LastSaveError string                           `json:"last_save_error"`
}

type StorageLocalDetailsResponseModel struct {
	Type string `json:"type"`
}

func marshallStorageLocalResourceModel(data StorageLocalResourceModel) map[string]any {
	return map[string]any{
		"name":        data.Name.ValueString(),
		"workspaceId": data.WorkspaceId.ValueString(),
		"settings":    map[string]any{"type": "LOCAL"},
	}
}

func (c *DatabasusClient) CreateStorageLocal(ctx context.Context, data StorageLocalResourceModel) (*StorageLocalResponseModel, error) {
	var result StorageLocalResponseModel
	b, _ := json.Marshal(marshallStorageLocalResourceModel(data))
	if err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetStorageLocal(ctx context.Context, id string) (*StorageLocalResponseModel, *ErrorDetails) {
	var result StorageLocalResponseModel
	if err := c.doRequest(ctx, "GET", "/storages/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateStorageLocal(ctx context.Context, id string, data StorageLocalResourceModel) (*StorageLocalResponseModel, error) {
	var result StorageLocalResponseModel
	body := marshallStorageLocalResourceModel(data)
	delete(body, "workspaceId")
	b, _ := json.Marshal(body)
	if err := c.doRequest(ctx, "PUT", "/storages/"+id, bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteStorageLocal(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/storages/"+id, nil, nil)
}
