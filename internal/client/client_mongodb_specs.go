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
*	MongoDB Database CRUD functions and model definitions
********************************************************/

type DatabaseMongoDbResourceModel struct {
	// database generic fields
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// DatabaseId     types.String `tfsdk:"database_id"`
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
	data.AuthDatabase = types.StringValue(response.MongoDb.AuthDatabase)
	data.Database = types.StringValue(response.MongoDb.Database)
	data.Host = types.StringValue(response.MongoDb.Host)
	data.IsHttps = types.BoolValue(response.MongoDb.IsHttps)
	data.Port = types.Int32Value(response.MongoDb.Port)
	data.Srv = types.BoolValue(response.MongoDb.IsSrv)
	data.DirectConnection = types.BoolValue(response.MongoDb.IsDirectConnection)
	data.CpuCount = types.Int32Value(response.MongoDb.CpuCount)
	// username and password are encrypted by databasus, we just ignore those for now since we could not detect changes anyways
}

type DatabaseMongoDbResponseModel struct {
	Id          string                              `json:"id"`
	Name        string                              `json:"name"`
	WorkspaceId string                              `json:"workspaceId"`
	Type        string                              `json:"type"`
	MongoDb     DatabaseMongoDbDetailsResponseModel `json:"mongodb"`
}

type DatabaseMongoDbDetailsResponseModel struct {
	Database           string `json:"database"`
	AuthDatabase       string `json:"authDatabase"`
	DatabaseId         string `json:"databaseId"`
	Host               string `json:"host"`
	Id                 string `json:"id"`
	IsSrv              bool   `json:"isSrv"`
	IsDirectConnection bool   `json:"isDirectConnection"`
	IsHttps            bool   `json:"isHttps"`
	Password           string `json:"password"`
	Port               int32  `json:"port"`
	Username           string `json:"username"`
	CpuCount           int32  `json:"cpuCount"`
	Version            string `json:"version"`
}

// internal helper to transform Model to map that can be used in request body.
func marshallDatabaseMongoDbResourceModel(data DatabaseMongoDbResourceModel) map[string]any {

	body := map[string]any{
		"isAgentTokenGenerated": true,
		"name":                  data.Name.ValueString(),
		"type":                  "MONGODB",
		"workspaceId":           data.WorkspaceId.ValueString(),
		"mongodb": map[string]any{
			"database":           data.Database.ValueString(),
			"cpuCount":           data.CpuCount.ValueInt32(),
			"authDatabase":       data.AuthDatabase.ValueString(),
			"host":               data.Host.ValueString(),
			"port":               data.Port.ValueInt32(),
			"username":           data.Username.ValueString(),
			"password":           data.Password.ValueString(),
			"isHttps":            data.IsHttps.ValueBool(),
			"isSrv":              data.Srv.ValueBool(),
			"isDirectConnection": data.DirectConnection.ValueBool(),
		},
	}
	return body
}

func (c *DatabasusClient) CreateDatabaseMongoDb(ctx context.Context, data DatabaseMongoDbResourceModel) (*DatabaseMongoDbResponseModel, error) {
	var result DatabaseMongoDbResponseModel
	body := marshallDatabaseMongoDbResourceModel(data)

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/create", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) GetDatabaseMongoDb(ctx context.Context, id string) (resultModel *DatabaseMongoDbResponseModel, errorD *ErrorDetails) {
	var result DatabaseMongoDbResponseModel

	err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateDatabaseMongoDb(ctx context.Context, id string, data DatabaseMongoDbResourceModel) (*DatabaseMongoDbResponseModel, error) {
	var result DatabaseMongoDbResponseModel

	body := marshallDatabaseMongoDbResourceModel(data)
	body["id"] = id

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "POST", "/databases/update", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete the configuration.
func (c *DatabasusClient) DeleteDatabaseMongoDb(ctx context.Context, id string) error {
	err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
