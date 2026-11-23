// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*******************************************************
* MongoDB Database CRUD functions and model definitions
********************************************************/

type DatabaseMongoDbResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	WorkspaceId      types.String `tfsdk:"workspace_id"`
	AuthDatabase     types.String `tfsdk:"auth_database"`
	CpuCount         types.Int32  `tfsdk:"cpu_count"`
	Database         types.String `tfsdk:"database"`
	Host             types.String `tfsdk:"host"`
	DirectConnection types.Bool   `tfsdk:"is_direct_connection"`
	Srv              types.Bool   `tfsdk:"is_srv"`
	IsHttps          types.Bool   `tfsdk:"is_https"`
	Port             types.Int32  `tfsdk:"port"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
}

func MapResponseToDatabaseMongoDbResourceModel(response *DatabaseMongoDbResponseModel, data *DatabaseMongoDbResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.AuthDatabase = types.StringValue(response.Settings.AuthDatabase)
	data.Database = types.StringValue(response.Settings.Database)
	data.Host = types.StringValue(response.Settings.Host)
	data.IsHttps = types.BoolValue(response.Settings.IsHttps)
	data.Port = types.Int32Value(response.Settings.Port)
	data.Srv = types.BoolValue(response.Settings.IsSrv)
	data.DirectConnection = types.BoolValue(response.Settings.IsDirectConnection)
	data.CpuCount = types.Int32Value(response.Settings.CpuCount)
}

type DatabaseMongoDbResponseModel struct {
	Id          string                           `json:"id"`
	Name        string                           `json:"name"`
	WorkspaceId string                           `json:"workspaceId"`
	Settings    DatabaseMongoDbDetailsResponseModel `json:"settings"`
	NotifierIds []string                         `json:"notifierIds"`
}

type DatabaseMongoDbDetailsResponseModel struct {
	Database           string `json:"database"`
	AuthDatabase       string `json:"authDatabase"`
	Host               string `json:"host"`
	Port               int32  `json:"port"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	IsSrv              bool   `json:"isSrv"`
	IsDirectConnection bool   `json:"isDirectConnection"`
	IsHttps            bool   `json:"isHttps"`
	CpuCount           int32  `json:"cpuCount"`
}

func marshallDatabaseMongoDbResourceModel(data DatabaseMongoDbResourceModel) map[string]any {
	return map[string]any{
		"workspaceId": data.WorkspaceId.ValueString(),
		"name":        data.Name.ValueString(),
		"notifierIds": []string{},
		"settings": map[string]any{
			"type":              "MONGODB",
			"host":              data.Host.ValueString(),
			"port":              data.Port.ValueInt32(),
			"username":          data.Username.ValueString(),
			"password":          data.Password.ValueString(),
			"database":          data.Database.ValueString(),
			"authDatabase":      data.AuthDatabase.ValueString(),
			"isHttps":           data.IsHttps.ValueBool(),
			"isSrv":             data.Srv.ValueBool(),
			"isDirectConnection": data.DirectConnection.ValueBool(),
			"cpuCount":          data.CpuCount.ValueInt32(),
			"sshTunnel":         map[string]any{"isEnabled": false},
			"excludeCollections": []string{},
		},
	}
}

func (c *DatabasusClient) CreateDatabaseMongoDb(ctx context.Context, data DatabaseMongoDbResourceModel) (*DatabaseMongoDbResponseModel, error) {
	var result DatabaseMongoDbResponseModel
	body := marshallDatabaseMongoDbResourceModel(data)
	b, _ := json.Marshal(body)
	if err := c.doRequest(ctx, "POST", "/databases", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetDatabaseMongoDb(ctx context.Context, id string) (*DatabaseMongoDbResponseModel, *ErrorDetails) {
	var result DatabaseMongoDbResponseModel
	if err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateDatabaseMongoDb(ctx context.Context, id string, data DatabaseMongoDbResourceModel) (*DatabaseMongoDbResponseModel, error) {
	var result DatabaseMongoDbResponseModel
	body := marshallDatabaseMongoDbResourceModel(data)
	delete(body, "workspaceId")
	b, _ := json.Marshal(map[string]any{"name": body["name"], "settings": body["settings"], "notifierIds": []string{}})
	if err := c.doRequest(ctx, "PUT", "/databases/"+id, bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteDatabaseMongoDb(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil)
}
