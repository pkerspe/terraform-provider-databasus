// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************
* PostgreSQL Database CRUD functions and model definitions
************************************************************/

type DatabasePostgresqlResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	WorkspaceId    types.String `tfsdk:"workspace_id"`
	Database       types.String `tfsdk:"database"`
	Host           types.String `tfsdk:"host"`
	IncludeSchemas types.List   `tfsdk:"include_schemas"`
	SslMode        types.String `tfsdk:"ssl_mode"`
	SslClientCert  types.String `tfsdk:"ssl_client_cert"`
	SslClientKey   types.String `tfsdk:"ssl_client_key"`
	SslRootCert    types.String `tfsdk:"ssl_root_cert"`
	Password       types.String `tfsdk:"password"`
	Port           types.Int32  `tfsdk:"port"`
	Username       types.String `tfsdk:"username"`
}

func MapResponseToDatabasePostgresqlResourceModel(response *DatabasePostgresqlResponseModel, data *DatabasePostgresqlResourceModel) {
	data.Id = types.StringValue(response.Id)
	data.Name = types.StringValue(response.Name)
	data.WorkspaceId = types.StringValue(response.WorkspaceId)
	data.Database = types.StringValue(response.Settings.Database)
	data.Host = types.StringValue(response.Settings.Host)

	if response.Settings.SslMode != "" {
		data.SslMode = types.StringValue(response.Settings.SslMode)
	}
	if response.Settings.SslClientCert != "" {
		data.SslClientCert = types.StringValue(response.Settings.SslClientCert)
	} else {
		data.SslClientCert = types.StringNull()
	}
	if response.Settings.SslClientKey != "" {
		data.SslClientKey = types.StringValue(response.Settings.SslClientKey)
	} else {
		data.SslClientKey = types.StringNull()
	}
	if response.Settings.SslRootCert != "" {
		data.SslRootCert = types.StringValue(response.Settings.SslRootCert)
	} else {
		data.SslRootCert = types.StringNull()
	}

	data.Port = types.Int32Value(int32(response.Settings.Port))
}

type DatabasePostgresqlResponseModel struct {
	Id          string                                 `json:"id"`
	Name        string                                 `json:"name"`
	WorkspaceId string                                 `json:"workspaceId"`
	Settings    DatabasePostgresqlDetailsResponseModel `json:"settings"`
	NotifierIds []string                               `json:"notifierIds"`
}

type DatabasePostgresqlDetailsResponseModel struct {
	Type               string   `json:"type"`
	Database           string   `json:"database"`
	Host               string   `json:"host"`
	IncludeSchemas     []string `json:"includeSchemas"`
	SslMode            string   `json:"sslMode"`
	SslClientCert      string   `json:"sslClientCert"`
	SslClientKey       string   `json:"sslClientKey"`
	SslRootCert        string   `json:"sslRootCert"`
	Password           string   `json:"password"`
	Port               int      `json:"port"`
	Username           string   `json:"username"`
	ExcludeTables      []string `json:"excludeTables"`
	ExcludeTableData   []string `json:"excludeTableData"`
	CpuCount           int      `json:"cpuCount"`
	IsSkipUserMappings bool     `json:"isSkipUserMappings"`
}

func marshallDatabasePostgresqlResourceModel(data DatabasePostgresqlResourceModel) map[string]any {
	includeSchemas := data.IncludeSchemas.Elements()
	includeSchemasStrings := make([]string, 0, len(includeSchemas))
	for _, item := range includeSchemas {
		strValue, ok := item.(types.String)
		if !ok {
			continue
		}
		includeSchemasStrings = append(includeSchemasStrings, strValue.ValueString())
	}

	return map[string]any{
		"workspaceId": data.WorkspaceId.ValueString(),
		"name":        data.Name.ValueString(),
		"notifierIds": []string{},
		"settings": map[string]any{
			"type":               "POSTGRES_LOGICAL",
			"host":               data.Host.ValueString(),
			"port":               data.Port.ValueInt32(),
			"username":           data.Username.ValueString(),
			"password":           data.Password.ValueString(),
			"database":           data.Database.ValueString(),
			"sslMode":            data.SslMode.ValueString(),
			"sslClientCert":      data.SslClientCert.ValueString(),
			"sslClientKey":       data.SslClientKey.ValueString(),
			"sslRootCert":        data.SslRootCert.ValueString(),
			"sshTunnel":          map[string]any{"isEnabled": false},
			"includeSchemas":     includeSchemasStrings,
			"excludeTables":      []string{},
			"excludeTableData":   []string{},
			"cpuCount":           1,
			"isSkipUserMappings": false,
		},
	}
}

func (c *DatabasusClient) CreateDatabasePostgresql(ctx context.Context, data DatabasePostgresqlResourceModel) (*DatabasePostgresqlResponseModel, error) {
	var result DatabasePostgresqlResponseModel
	body := marshallDatabasePostgresqlResourceModel(data)
	b, _ := json.Marshal(body)
	if err := c.doRequest(ctx, "POST", "/databases", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetDatabasePostgresql(ctx context.Context, id string) (resultModel *DatabasePostgresqlResponseModel, errorD *ErrorDetails) {
	var result DatabasePostgresqlResponseModel
	if err := c.doRequest(ctx, "GET", "/databases/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateDatabasePostgresql(ctx context.Context, id string, data DatabasePostgresqlResourceModel) (*DatabasePostgresqlResponseModel, error) {
	var result DatabasePostgresqlResponseModel
	body := marshallDatabasePostgresqlResourceModel(data)
	delete(body, "workspaceId")
	b, _ := json.Marshal(map[string]any{
		"name":        body["name"],
		"settings":    body["settings"],
		"notifierIds": []string{},
	})
	if err := c.doRequest(ctx, "PUT", "/databases/"+id, bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteDatabasePostgresql(ctx context.Context, id string) error {
	if err := c.doRequest(ctx, "DELETE", "/databases/"+id, nil, nil); err != nil {
		return err
	}
	return nil
}
