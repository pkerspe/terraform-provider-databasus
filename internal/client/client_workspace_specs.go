// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkspaceResponseModel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type WorkspacesListResponse struct {
	Items      []WorkspaceResponseModel `json:"items"`
	NextCursor *string                  `json:"nextCursor"`
}

type WorkspaceDataSourceModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
}

func (c *DatabasusClient) CreateWorkspace(ctx context.Context, name string) (*WorkspaceResponseModel, error) {
	var result WorkspaceResponseModel
	b, _ := json.Marshal(map[string]string{"name": name})
	if err := c.doRequest(ctx, "POST", "/workspaces", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetWorkspace(ctx context.Context, id string) (*WorkspaceResponseModel, *ErrorDetails) {
	var result WorkspaceResponseModel
	if err := c.doRequest(ctx, "GET", "/workspaces/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) GetWorkspaceByName(ctx context.Context, name string) (*WorkspaceResponseModel, *ErrorDetails) {
	var cursor *string
	for {
		path := "/workspaces?limit=100"
		if cursor != nil && *cursor != "" {
			path += "&after=" + url.QueryEscape(*cursor)
		}
		var page WorkspacesListResponse
		if err := c.doRequest(ctx, "GET", path, nil, &page); err != nil {
			return nil, err
		}
		for _, workspace := range page.Items {
			if workspace.Name == name {
				return &workspace, nil
			}
		}
		if page.NextCursor == nil || *page.NextCursor == "" {
			break
		}
		cursor = page.NextCursor
	}
	return nil, &ErrorDetails{ErrorCode: 404, ResponseBody: "Workspace not found for name: " + name, ErrorInst: fmt.Errorf("workspace not found for name: %s", name)}
}

func (c *DatabasusClient) ListWorkspaces(ctx context.Context) (*WorkspacesListResponse, error) {
	var result WorkspacesListResponse
	var items []WorkspaceResponseModel
	var cursor *string
	for {
		path := "/workspaces?limit=100"
		if cursor != nil && *cursor != "" {
			path += "&after=" + url.QueryEscape(*cursor)
		}
		var page WorkspacesListResponse
		if err := c.doRequest(ctx, "GET", path, nil, &page); err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if page.NextCursor == nil || *page.NextCursor == "" {
			break
		}
		cursor = page.NextCursor
	}
	result.Items = items
	return &result, nil
}

func (c *DatabasusClient) UpdateWorkspace(ctx context.Context, id string, name string) (*WorkspaceResponseModel, error) {
	var result WorkspaceResponseModel
	b, _ := json.Marshal(map[string]string{"name": name})
	if err := c.doRequest(ctx, "PUT", "/workspaces/"+id, bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteWorkspace(ctx context.Context, id string) error {
	if err := c.doRequest(ctx, "DELETE", "/workspaces/"+id, nil, nil); err != nil {
		return err
	}
	return nil
}
