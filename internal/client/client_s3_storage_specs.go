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
*								S3 Storage CRUD functions
********************************************************/

type StorageS3ResourceModel struct {
	// storage generic fields
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	IsSystem      types.Bool   `tfsdk:"is_system"`
	LastSaveError types.String `tfsdk:"last_save_error"`
	// storage specific fields
	S3AccessKey             types.String `tfsdk:"s3_access_key"`
	S3Bucket                types.String `tfsdk:"s3_bucket"`
	S3Endpoint              types.String `tfsdk:"s3_endpoint"`
	S3Prefix                types.String `tfsdk:"s3_prefix"`
	S3Region                types.String `tfsdk:"s3_region"`
	S3SecretKey             types.String `tfsdk:"s3_secret_key"`
	S3StorageClass          types.String `tfsdk:"s3_storage_class"`
	S3UseVirtualHostedStyle types.Bool   `tfsdk:"s3_use_virtual_hosted_style"`
	SkipTLSVerify           types.Bool   `tfsdk:"skip_tls_verify"`
	StorageId               types.String `tfsdk:"storage_id"`
}

func MapResponseToStorageS3ResourceModel(response *StorageS3ResponseModel, data *StorageS3ResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.IsSystem = types.BoolValue(response.IsSystem)
	data.LastSaveError = types.StringValue(response.LastSaveError)
	data.S3Bucket = types.StringValue(response.S3Storage.S3Bucket)
	data.S3Endpoint = types.StringValue(response.S3Storage.S3Endpoint)
	data.S3Prefix = types.StringValue(response.S3Storage.S3Prefix)
	data.S3Region = types.StringValue(response.S3Storage.S3Region)
	data.S3StorageClass = types.StringValue(response.S3Storage.S3StorageClass)
	data.S3UseVirtualHostedStyle = types.BoolValue(response.S3Storage.S3UseVirtualHostedStyle)
	data.SkipTLSVerify = types.BoolValue(response.S3Storage.SkipTLSVerify)
	data.StorageId = types.StringValue(response.S3Storage.StorageId)
	// The following two are internally encoded by databasus, so we do not store to the state since it would cause a mismatch from the planned values
	// this could be refactored in the future to store the encoded values in a separate property, which is computed, but do we really need it in the state?
	// data.S3AccessKey = types.StringValue(response.S3Storage.S3AccessKey)
	// data.S3SecretKey = types.StringValue(response.S3Storage.S3SecretKey)
}

type StorageS3ResponseModel struct {
	Id            string                        `json:"id"`
	Name          string                        `json:"name"`
	WorkspaceId   string                        `json:"workspaceId"`
	IsSystem      bool                          `json:"isSystem"`
	S3Storage     StorageS3DetailsResponseModel `json:"s3Storage"`
	LastSaveError string                        `json:"lastSaveError"`
}

type StorageS3DetailsResponseModel struct {
	S3AccessKey             string `json:"s3AccessKey"`
	S3Bucket                string `json:"s3Bucket"`
	S3Endpoint              string `json:"s3Endpoint"`
	S3Prefix                string `json:"s3Prefix"`
	S3Region                string `json:"s3Region"`
	S3SecretKey             string `json:"s3SecretKey"`
	S3StorageClass          string `json:"s3StorageClass"`
	S3UseVirtualHostedStyle bool   `json:"s3UseVirtualHostedStyle"`
	SkipTLSVerify           bool   `json:"skipTLSVerify"`
	StorageId               string `json:"storageId"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallStorageS3ResourceModel(data StorageS3ResourceModel) map[string]any {
	body := map[string]any{
		"isSystem":      data.IsSystem.ValueBool(),
		"lastSaveError": "",
		"name":          data.Name.ValueString(),
		"type":          "S3",
		"workspaceId":   data.WorkspaceId.ValueString(),
		"s3Storage": map[string]any{
			"s3AccessKey":             data.S3AccessKey.ValueString(),
			"s3Bucket":                data.S3Bucket.ValueString(),
			"s3Endpoint":              data.S3Endpoint.ValueString(),
			"s3Prefix":                data.S3Prefix.ValueString(),
			"s3Region":                data.S3Region.ValueString(),
			"s3SecretKey":             data.S3SecretKey.ValueString(),
			"s3StorageClass":          data.S3StorageClass.ValueString(),
			"s3UseVirtualHostedStyle": data.S3UseVirtualHostedStyle.ValueBool(),
			"skipTLSVerify":           data.SkipTLSVerify.ValueBool(),
			// "storageId":               data.StorageId.ValueString(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateStorageS3(ctx context.Context, data StorageS3ResourceModel) (*StorageS3ResponseModel, error) {
	var result StorageS3ResponseModel
	body := marshallStorageS3ResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetStorageS3(ctx context.Context, id string) (resultModel *StorageS3ResponseModel, errorD *ErrorDetails) {
	var result StorageS3ResponseModel

	err := c.doRequest(ctx, "GET", "/storages/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateStorageS3(ctx context.Context, id string, data StorageS3ResourceModel) (*StorageS3ResponseModel, error) {
	var result StorageS3ResponseModel

	body := marshallStorageS3ResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the S3 storage configuration.
func (c *DatabasusClient) DeleteStorageS3(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/storages/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
