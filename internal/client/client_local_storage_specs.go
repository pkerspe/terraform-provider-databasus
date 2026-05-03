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
*		Local Storage CRUD functions and model definitions
********************************************************/

type StorageLocalResourceModel struct {
	// storage generic fields
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	LastSaveError types.String `tfsdk:"last_save_error"`
	//StorageId     types.String `tfsdk:"storage_id"`
}

func MapResponseToStorageLocalResourceModel(response *StorageLocalResponseModel, data *StorageLocalResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.LastSaveError = types.StringValue(response.LastSaveError)
	//data.StorageId = types.StringValue(response.LocalStorage.StorageId)
}

type StorageLocalResponseModel struct {
	Id            string                           `json:"id"`
	Name          string                           `json:"name"`
	WorkspaceId   string                           `json:"workspaceId"`
	LocalStorage  StorageLocalDetailsResponseModel `json:"localStorage"`
	LastSaveError string                           `json:"lastSaveError"`
}

type StorageLocalDetailsResponseModel struct {
	StorageId string `json:"storageId"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallStorageLocalResourceModel(data StorageLocalResourceModel) map[string]any {
	body := map[string]any{
		"isSystem":      false,
		"lastSaveError": "",
		"name":          data.Name.ValueString(),
		"type":          "LOCAL",
		"workspaceId":   data.WorkspaceId.ValueString(),
		"localStorage":  map[string]any{
			// "storageId":               data.StorageId.ValueString(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateStorageLocal(ctx context.Context, data StorageLocalResourceModel) (*StorageLocalResponseModel, error) {
	var result StorageLocalResponseModel
	body := marshallStorageLocalResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetStorageLocal(ctx context.Context, id string) (resultModel *StorageLocalResponseModel, errorD *ErrorDetails) {
	var result StorageLocalResponseModel

	err := c.doRequest(ctx, "GET", "/storages/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateStorageLocal(ctx context.Context, id string, data StorageLocalResourceModel) (*StorageLocalResponseModel, error) {
	var result StorageLocalResponseModel

	body := marshallStorageLocalResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the Local storage configuration.
func (c *DatabasusClient) DeleteStorageLocal(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/storages/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
