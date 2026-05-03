// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************
*	PostgreSQL Database CRUD functions and model definitions
************************************************************/

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
