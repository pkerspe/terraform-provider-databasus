// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomRoundTripper is a custom implementation of http.RoundTripper
// that adds the Authorization header with the JWT token.
type CustomRoundTripper struct {
	Transport http.RoundTripper
	Token     string
}

// RoundTrip executes a single HTTP transaction and adds the Authorization header.
func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Use the original transport to execute the request
	return c.Transport.RoundTrip(req)
}

type DatabasusClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

type ErrorDetails struct {
	ErrorCode    int
	ResponseBody string
	ErrorInst    error
}

func (e *ErrorDetails) Error() string {
	return e.ErrorInst.Error()
}

func (e *ErrorDetails) IsNotFound() bool {
	return e.ErrorCode == 404 ||
		strings.Contains(strings.ToLower(e.ResponseBody), "record not found") //fix for invalid API RCs from Databasus
}

func NewDatabasusClient(baseURL, token string) *DatabasusClient {
	// Create a new HTTP client with the custom RoundTripper
	client := &http.Client{
		Transport: &CustomRoundTripper{
			Transport: http.DefaultTransport, // Use the default transport
			Token:     token,
		},
	}
	return &DatabasusClient{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    client,
	}
}

func GetJWT(baseURL, email, password string) (string, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	b, _ := json.Marshal(body)

	resp, err := http.Post(baseURL+"/users/signin", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token  string `json:"token"`
		UserId string `json:"userId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *DatabasusClient) doRequest(ctx context.Context, method, path string, body io.Reader, out interface{}) *ErrorDetails {
	var errorDetails ErrorDetails
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		errorDetails.ErrorInst = fmt.Errorf("failed to create request: %w", err)
		return &errorDetails
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		errorDetails.ErrorInst = fmt.Errorf("request failed: %w", err)
		return &errorDetails
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		errorDetails.ResponseBody = string(respBody)
		errorDetails.ErrorCode = resp.StatusCode
		errorDetails.ErrorInst = fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return &errorDetails
		//return fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBody)), &errorDetails
	}

	// Decode if output is provided
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			errorDetails.ErrorInst = fmt.Errorf("failed to decode response: %w", err)
			return &errorDetails
		}
	}

	return nil
}

/*******************************************************
*								WORKSPACE CRUD functions
********************************************************/

type WorkspaceResponseModel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type WorkspacesListResponse struct {
	Items []WorkspaceResponseModel `json:"workspaces"`
}

type WorkspaceDataSourceModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
}

// Create a new Workspace in Databasus for the given name.
func (c *DatabasusClient) CreateWorkspace(ctx context.Context, name string) (*WorkspaceResponseModel, error) {
	var result WorkspaceResponseModel

	body := map[string]string{
		"name": name,
	}

	b, _ := json.Marshal(body)

	err := c.doRequest(ctx, "POST", "/workspaces", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetWorkspace(ctx context.Context, id string) (*WorkspaceResponseModel, error) {
	var result WorkspaceResponseModel

	err := c.doRequest(ctx, "GET", "/workspaces/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) ListWorkspaces(ctx context.Context) (*WorkspacesListResponse, error) {
	var result WorkspacesListResponse

	err := c.doRequest(ctx, "GET", "/workspaces", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Update an Existing Workspace from Databasus for a given id.
func (c *DatabasusClient) UpdateWorkspace(ctx context.Context, id string, name string) (*WorkspaceResponseModel, error) {
	var result WorkspaceResponseModel

	body := map[string]string{
		"name": name,
	}

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "PUT", "/workspaces/"+id, bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete an Existing Workspace from Databasus for a given id.
func (c *DatabasusClient) DeleteWorkspace(ctx context.Context, id string) error {
	var result WorkspaceResponseModel

	err := c.doRequest(ctx, "DELETE", "/workspaces/"+id, nil, &result)
	if err != nil {
		return err
	}

	return nil
}

/*******************************************************
*								SETTINGS CRUD functions
********************************************************/

type SettingsResponseModel struct {
	Id                                string `json:"id"`
	IsAllowExternalRegistrations      bool   `json:"isAllowExternalRegistrations"`
	IsAllowMemberInvitations          bool   `json:"isAllowMemberInvitations"`
	IsMemberAllowedToCreateWorkspaces bool   `json:"isMemberAllowedToCreateWorkspaces"`
}

type SettingsDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
}

type SettingsResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
}

// settings always exist, so the create just internally calls the update function.
func (c *DatabasusClient) CreateUsersSettings(ctx context.Context, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool) (*SettingsResponseModel, error) {
	currentSettings, err := c.GetUsersSettings(ctx)
	if err != nil {
		return nil, err
	}

	return c.UpdateUsersSettings(ctx, currentSettings.Id, allowExternalRegistrations, allowMemberInvitations, memberAllowedToCreateWorkspaces)
}

func (c *DatabasusClient) GetUsersSettings(ctx context.Context) (*SettingsResponseModel, error) {
	var result SettingsResponseModel

	err := c.doRequest(ctx, "GET", "/users/settings", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateUsersSettings(ctx context.Context, id string, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool) (*SettingsResponseModel, error) {
	var result SettingsResponseModel

	body := map[string]any{
		"isAllowExternalRegistrations":      allowExternalRegistrations,
		"isAllowMemberInvitations":          allowMemberInvitations,
		"isMemberAllowedToCreateWorkspaces": memberAllowedToCreateWorkspaces,
	}

	if id != "" {
		body["id"] = id
	}

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "PUT", "/users/settings", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete does not exist, so we just do nothing here.
func (c *DatabasusClient) DeleteUsersSettings(ctx context.Context, id string) error {
	return nil
}

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

/*******************************************************
*								Local Storage CRUD functions
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
