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
*	 MySQL Database CRUD functions and model definitions
********************************************************/

type DatabaseMySqlResourceModel struct {
	// database generic fields
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// DatabaseId     types.String `tfsdk:"database_id"`
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
	data.Database = types.StringValue(response.MySql.Database)
	data.Host = types.StringValue(response.MySql.Host)
	data.IsHttps = types.BoolValue(response.MySql.IsHttps)
	data.Port = types.Int32Value(response.MySql.Port)
	// username and password are encrypted by databasus, we just ignore those for now since we could not detect changes anyways
}

type DatabaseMySqlResponseModel struct {
	Id          string                            `json:"id"`
	Name        string                            `json:"name"`
	WorkspaceId string                            `json:"workspaceId"`
	Type        string                            `json:"type"`
	MySql       DatabaseMySqlDetailsResponseModel `json:"mysql"`
}

type DatabaseMySqlDetailsResponseModel struct {
	Database   string `json:"database"`
	DatabaseId string `json:"databaseId"`
	Host       string `json:"host"`
	Id         string `json:"id"`
	IsHttps    bool   `json:"isHttps"`
	Password   string `json:"password"`
	Port       int32  `json:"port"`
	Privileges string `json:"privileges"`
	Username   string `json:"username"`
	Version    string `json:"version"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallDatabaseMySqlResourceModel(data DatabaseMySqlResourceModel) map[string]any {

	body := map[string]any{
		"isAgentTokenGenerated": true,
		"name":                  data.Name.ValueString(),
		"type":                  "MYSQL",
		"workspaceId":           data.WorkspaceId.ValueString(),
		"mysql": map[string]any{
			"database": data.Database.ValueString(),
			"host":     data.Host.ValueString(),
			"port":     data.Port.ValueInt32(),
			"username": data.Username.ValueString(),
			"password": data.Password.ValueString(),
			"isHttps":  data.IsHttps.ValueBool(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateDatabaseMySql(ctx context.Context, data DatabaseMySqlResourceModel) (*DatabaseMySqlResponseModel, error) {
	var result DatabaseMySqlResponseModel
	body := marshallDatabaseMySqlResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/create", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetDatabaseMySql(ctx context.Context, id string) (resultModel *DatabaseMySqlResponseModel, errorD *ErrorDetails) {
	var result DatabaseMySqlResponseModel

	err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateDatabaseMySql(ctx context.Context, id string, data DatabaseMySqlResourceModel) (*DatabaseMySqlResponseModel, error) {
	var result DatabaseMySqlResponseModel

	body := marshallDatabaseMySqlResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/update", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the configuration.
func (c *DatabasusClient) DeleteDatabaseMySql(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
