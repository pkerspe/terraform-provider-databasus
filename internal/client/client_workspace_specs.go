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
*		 WORKSPACE CRUD functions and model definitions
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

func (c *DatabasusClient) GetWorkspace(ctx context.Context, id string) (*WorkspaceResponseModel, *ErrorDetails) {
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
