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
* MySQL Database CRUD functions and model definitions
********************************************************/

type DatabaseMySqlResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	Database    types.String `tfsdk:"database"`
	Host        types.String `tfsdk:"host"`
	IsHttps     types.Bool   `tfsdk:"is_https"`
	Port        types.Int32  `tfsdk:"port"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
}

func MapResponseToDatabaseMySqlResourceModel(response *DatabaseMySqlResponseModel, data *DatabaseMySqlResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.Database = types.StringValue(response.Settings.Database)
	data.Host = types.StringValue(response.Settings.Host)
	data.IsHttps = types.BoolValue(response.Settings.IsHttps)
	data.Port = types.Int32Value(response.Settings.Port)
}

type DatabaseMySqlResponseModel struct {
	Id          string                            `json:"id"`
	Name        string                            `json:"name"`
	WorkspaceId string                            `json:"workspaceId"`
	Settings    DatabaseMySqlDetailsResponseModel `json:"settings"`
	NotifierIds []string                          `json:"notifierIds"`
}

type DatabaseMySqlDetailsResponseModel struct {
	Type     string `json:"type"`
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsHttps  bool   `json:"isHttps"`
}

func marshallDatabaseMySqlResourceModel(data DatabaseMySqlResourceModel) map[string]any {
	return map[string]any{
		"workspaceId": data.WorkspaceId.ValueString(),
		"name":        data.Name.ValueString(),
		"notifierIds": []string{},
		"settings": map[string]any{
			"type":          "MYSQL",
			"host":          data.Host.ValueString(),
			"port":          data.Port.ValueInt32(),
			"username":      data.Username.ValueString(),
			"password":      data.Password.ValueString(),
			"database":      data.Database.ValueString(),
			"isHttps":       data.IsHttps.ValueBool(),
			"sshTunnel":     map[string]any{"isEnabled": false},
			"excludeTables": []string{},
		},
	}
}

func (c *DatabasusClient) CreateDatabaseMySql(ctx context.Context, data DatabaseMySqlResourceModel) (*DatabaseMySqlResponseModel, error) {
	var result DatabaseMySqlResponseModel
	body := marshallDatabaseMySqlResourceModel(data)
	b, _ := json.Marshal(body)
	if err := c.doRequest(ctx, "POST", "/databases", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetDatabaseMySql(ctx context.Context, id string) (*DatabaseMySqlResponseModel, *ErrorDetails) {
	var result DatabaseMySqlResponseModel
	if err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateDatabaseMySql(ctx context.Context, id string, data DatabaseMySqlResourceModel) (*DatabaseMySqlResponseModel, error) {
	var result DatabaseMySqlResponseModel
	body := marshallDatabaseMySqlResourceModel(data)
	delete(body, "workspaceId")
	b, _ := json.Marshal(map[string]any{"name": body["name"], "settings": body["settings"], "notifierIds": []string{}})
	if err := c.doRequest(ctx, "PUT", "/databases/"+id, bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteDatabaseMySql(ctx context.Context, id string) error {
	if err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil); err != nil {
		return err
	}
	return nil
}
