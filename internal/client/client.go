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

/*******************************************************
*	   			PostgreSQL Database CRUD functions
********************************************************/

type DatabasePostgresqlResourceModel struct {
	// database generic fields
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// Type           types.String `tfsdk:"type"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	Database    types.String `tfsdk:"database"`
	// DatabaseId     types.String `tfsdk:"database_id"`
	Host           types.String `tfsdk:"host"`
	IncludeSchemas types.List   `tfsdk:"include_schemas"`
	IsHttps        types.Bool   `tfsdk:"is_https"`
	Password       types.String `tfsdk:"password"`
	Port           types.Int32  `tfsdk:"port"`
	Username       types.String `tfsdk:"username"`
}

func MapResponseToDatabasePostgresqlResourceModel(response *DatabasePostgresqlResponseModel, data *DatabasePostgresqlResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.Database = types.StringValue(response.Postgresql.Database)
	// data.Type = types.StringValue(response.Type)
	data.Host = types.StringValue(response.Postgresql.Host)
	data.IsHttps = types.BoolValue(response.Postgresql.IsHttps)
	data.Port = types.Int32Value(int32(response.Postgresql.Port))
	// TODO: CHECK how to map: data.IncludeSchemas = types.ListValue(types.String(), response.Postgresql.IncludeSchemas)
	// username and password are encrypted by databasus, we just ignore those for now since we could not detect changes anyways
}

type DatabasePostgresqlResponseModel struct {
	Id          string                                 `json:"id"`
	Name        string                                 `json:"name"`
	WorkspaceId string                                 `json:"workspaceId"`
	Type        string                                 `json:"type"`
	Postgresql  DatabasePostgresqlDetailsResponseModel `json:"postgresql"`
}

type DatabasePostgresqlDetailsResponseModel struct {
	Database       string   `json:"database"`
	DatabaseId     string   `json:"databaseId"`
	Host           string   `json:"host"`
	Id             string   `json:"id"`
	IncludeSchemas []string `json:"includeSchemas"`
	IsHttps        bool     `json:"isHttps"`
	Password       string   `json:"password"`
	Port           int      `json:"port"`
	Username       string   `json:"username"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallDatabasePostgresqlResourceModel(data DatabasePostgresqlResourceModel) map[string]any {
	// Extract the elements as a []attr.Value (which can be any value from the Terraform SDK)
	includeSchemas := data.IncludeSchemas.Elements()

	// Convert []attr.Value to []string
	var includeSchemasStrings []string
	for _, item := range includeSchemas {
		// Use the ValueString() method, which works for attributes like types.String
		strValue, ok := item.(types.String)
		if !ok {
			// If the element is not a types.String, handle this case
			return nil // Or handle the error accordingly
		}
		includeSchemasStrings = append(includeSchemasStrings, strValue.ValueString())
	}

	body := map[string]any{
		"isAgentTokenGenerated": true,
		"name":                  data.Name.ValueString(),
		"type":                  "POSTGRES",
		"workspaceId":           data.WorkspaceId.ValueString(),
		// "notifiers": []
		"postgresql": map[string]any{
			"backupType":          "PG_DUMP",
			"cpuCount":            1,
			"database":            data.Database.ValueString(),
			"host":                data.Host.ValueString(),
			"port":                data.Port.ValueInt32(),
			"username":            data.Username.ValueString(),
			"password":            data.Password.ValueString(),
			"isHttps":             data.IsHttps.ValueBool(),
			"isExcludeExtensions": true,
			"includeSchemas":      includeSchemasStrings,
			//"databaseId"
			//"version"
		},
	}
	return body
}

func (c *DatabasusClient) CreateDatabasePostgresql(ctx context.Context, data DatabasePostgresqlResourceModel) (*DatabasePostgresqlResponseModel, error) {
	var result DatabasePostgresqlResponseModel
	body := marshallDatabasePostgresqlResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/create", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetDatabasePostgresql(ctx context.Context, id string) (resultModel *DatabasePostgresqlResponseModel, errorD *ErrorDetails) {
	var result DatabasePostgresqlResponseModel

	err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateDatabasePostgresql(ctx context.Context, id string, data DatabasePostgresqlResourceModel) (*DatabasePostgresqlResponseModel, error) {
	var result DatabasePostgresqlResponseModel

	body := marshallDatabasePostgresqlResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/update", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the configuration.
func (c *DatabasusClient) DeleteDatabasePostgresql(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

/*******************************************************
*	   			Webhook Notifier CRUD functions
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

/*******************************************************
*								Backup Config CRUD functions
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

// Delete is not possible for backup configs, can only set to false, but need to use update then
func (c *DatabasusClient) DeleteBackupConfig(ctx context.Context) error {
	return nil
}
