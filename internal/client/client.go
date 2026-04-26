// Copyright (c) KerspeP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DatabasusClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

type WorkspaceResponseModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type WorkspacesListResponse struct {
	Items []WorkspaceResponseModel `json:"workspaces"`
}

type WorkspaceDataSourceModel struct {
	CreatedAt types.String `tfsdk:"createdAt"`
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
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

func (c *DatabasusClient) doRequest(ctx context.Context, method, path string, body io.Reader, out interface{}) error {
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	// Decode if output is provided
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
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
