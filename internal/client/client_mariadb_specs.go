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
*	MariaDB Database CRUD functions and model definitions
********************************************************/

type DatabaseMariaDbResourceModel struct {
	// database generic fields
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// DatabaseId     types.String `tfsdk:"database_id"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	Database      types.String `tfsdk:"database"`
	Host          types.String `tfsdk:"host"`
	IsHttps       types.Bool   `tfsdk:"is_https"`
	Port          types.Int32  `tfsdk:"port"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	ExcludeEvents types.Bool   `tfsdk:"exclude_events"`
}

func MapResponseToDatabaseMariaDbResourceModel(response *DatabaseMariaDbResponseModel, data *DatabaseMariaDbResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.Database = types.StringValue(response.MariaDb.Database)
	data.Host = types.StringValue(response.MariaDb.Host)
	data.IsHttps = types.BoolValue(response.MariaDb.IsHttps)
	data.Port = types.Int32Value(response.MariaDb.Port)
	data.ExcludeEvents = types.BoolValue(response.MariaDb.IsExcludeEvents)
	// TODO: CHECK how to map: data.IncludeSchemas = types.ListValue(types.String(), response.MariaDb.IncludeSchemas)
	// username and password are encrypted by databasus, we just ignore those for now since we could not detect changes anyways
}

type DatabaseMariaDbResponseModel struct {
	Id          string                              `json:"id"`
	Name        string                              `json:"name"`
	WorkspaceId string                              `json:"workspaceId"`
	Type        string                              `json:"type"`
	MariaDb     DatabaseMariaDbDetailsResponseModel `json:"mariaDb"`
}

type DatabaseMariaDbDetailsResponseModel struct {
	Database        string `json:"database"`
	DatabaseId      string `json:"databaseId"`
	Host            string `json:"host"`
	Id              string `json:"id"`
	IsExcludeEvents bool   `json:"isExcludeEvents"`
	IsHttps         bool   `json:"isHttps"`
	Password        string `json:"password"`
	Port            int32  `json:"port"`
	Username        string `json:"username"`
	Privileges      string `json:"privileges"`
	Version         string `json:"version"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallDatabaseMariaDbResourceModel(data DatabaseMariaDbResourceModel) map[string]any {

	body := map[string]any{
		"isAgentTokenGenerated": true,
		"name":                  data.Name.ValueString(),
		"type":                  "MARIADB",
		"workspaceId":           data.WorkspaceId.ValueString(),
		"mariaDb": map[string]any{
			"database":        data.Database.ValueString(),
			"host":            data.Host.ValueString(),
			"port":            data.Port.ValueInt32(),
			"username":        data.Username.ValueString(),
			"password":        data.Password.ValueString(),
			"isHttps":         data.IsHttps.ValueBool(),
			"isExcludeEvents": data.ExcludeEvents.ValueBool(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateDatabaseMariaDb(ctx context.Context, data DatabaseMariaDbResourceModel) (*DatabaseMariaDbResponseModel, error) {
	var result DatabaseMariaDbResponseModel
	body := marshallDatabaseMariaDbResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/create", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetDatabaseMariaDb(ctx context.Context, id string) (resultModel *DatabaseMariaDbResponseModel, errorD *ErrorDetails) {
	var result DatabaseMariaDbResponseModel

	err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateDatabaseMariaDb(ctx context.Context, id string, data DatabaseMariaDbResourceModel) (*DatabaseMariaDbResponseModel, error) {
	var result DatabaseMariaDbResponseModel

	body := marshallDatabaseMariaDbResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/update", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the configuration.
func (c *DatabasusClient) DeleteDatabaseMariaDb(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
