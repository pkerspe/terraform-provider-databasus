// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageS3ResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	WorkspaceId             types.String `tfsdk:"workspace_id"`
	LastSaveError           types.String `tfsdk:"last_save_error"`
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
	// The current API identifies the storage by the top-level ID; the legacy computed attribute mirrors it.
	data.StorageId = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.LastSaveError = types.StringValue(response.LastSaveError)
	data.S3Bucket = types.StringValue(response.Settings.S3Bucket)
	data.S3Endpoint = types.StringValue(response.Settings.S3Endpoint)
	data.S3Prefix = types.StringValue(response.Settings.S3Prefix)
	data.S3Region = types.StringValue(response.Settings.S3Region)
	data.S3StorageClass = types.StringValue(response.Settings.S3StorageClass)
	data.S3UseVirtualHostedStyle = types.BoolValue(response.Settings.IsVirtualHostedStyle)
	data.SkipTLSVerify = types.BoolValue(response.Settings.IsSkipTlsVerify)
}

type StorageS3ResponseModel struct {
	Id            string                        `json:"id"`
	Name          string                        `json:"name"`
	WorkspaceId   string                        `json:"workspaceId"`
	Settings      StorageS3DetailsResponseModel `json:"settings"`
	LastSaveError string                        `json:"lastSaveError"`
}

type StorageS3DetailsResponseModel struct {
	Type                 string `json:"type"`
	S3AccessKey          string `json:"s3AccessKey"`
	S3Bucket             string `json:"s3Bucket"`
	S3Endpoint           string `json:"s3Endpoint"`
	S3Prefix             string `json:"s3Prefix"`
	S3Region             string `json:"s3Region"`
	S3SecretKey          string `json:"s3SecretKey"`
	S3StorageClass       string `json:"s3StorageClass"`
	IsVirtualHostedStyle bool   `json:"isVirtualHostedStyle"`
	IsSkipTlsVerify      bool   `json:"isSkipTlsVerify"`
}

func marshallStorageS3ResourceModel(data StorageS3ResourceModel) map[string]any {
	settings := map[string]any{
		"type":                 "S3",
		"s3AccessKey":          data.S3AccessKey.ValueString(),
		"s3Bucket":             data.S3Bucket.ValueString(),
		"s3Endpoint":           data.S3Endpoint.ValueString(),
		"s3Prefix":             data.S3Prefix.ValueString(),
		"s3Region":             data.S3Region.ValueString(),
		"s3SecretKey":          data.S3SecretKey.ValueString(),
		"s3StorageClass":       data.S3StorageClass.ValueString(),
		"isVirtualHostedStyle": data.S3UseVirtualHostedStyle.ValueBool(),
		"isSkipTlsVerify":      data.SkipTLSVerify.ValueBool(),
	}
	return map[string]any{"name": data.Name.ValueString(), "workspaceId": data.WorkspaceId.ValueString(), "settings": settings}
}

func (c *DatabasusClient) CreateStorageS3(ctx context.Context, data StorageS3ResourceModel) (*StorageS3ResponseModel, error) {
	var result StorageS3ResponseModel
	body, err := json.Marshal(marshallStorageS3ResourceModel(data))
	if err != nil {
		return nil, err
	}
	if err := c.doRequest(ctx, "POST", "/storages", bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetStorageS3(ctx context.Context, id string) (*StorageS3ResponseModel, *ErrorDetails) {
	var result StorageS3ResponseModel
	if err := c.doRequest(ctx, "GET", "/storages/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateStorageS3(ctx context.Context, id string, data StorageS3ResourceModel) (*StorageS3ResponseModel, error) {
	var result StorageS3ResponseModel
	body := marshallStorageS3ResourceModel(data)
	delete(body, "workspaceId")
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if err := c.doRequest(ctx, "PUT", "/storages/"+id, bytes.NewBuffer(encodedBody), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteStorageS3(ctx context.Context, id string) *ErrorDetails {
	return c.doRequest(ctx, "DELETE", "/storages/"+id, nil, nil)
}
