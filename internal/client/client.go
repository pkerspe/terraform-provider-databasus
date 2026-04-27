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

// CustomRoundTripper is a custom implementation of http.RoundTripper
// that adds the Authorization header with the JWT token.
type CustomRoundTripper struct {
	Transport http.RoundTripper
	Token     string
}

// RoundTrip executes a single HTTP transaction and adds the Authorization header.
func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Use the original transport to execute the request
	return c.Transport.RoundTrip(req)
}

type DatabasusClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewDatabasusClient(baseURL, token string) *DatabasusClient {
	// Create a new HTTP client with the custom RoundTripper
	client := &http.Client{
		Transport: &CustomRoundTripper{
			Transport: http.DefaultTransport, // Use the default transport
			Token:     token,
		},
	}
	return &DatabasusClient{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    client,
	}
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
	CreatedAt types.String `tfsdk:"created_at"`
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

// Create a new Workspace in Databasus for the given name
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

// Delete an Existing Workspace from Databasus for a given id
func (c *DatabasusClient) DeleteWorkspace(ctx context.Context, id string) error {
	var result WorkspaceResponseModel

	err := c.doRequest(ctx, "DELETE", "/workspaces/"+id, nil, &result)
	if err != nil {
		return err
	}

	return nil
}

// Update an Existing Workspace from Databasus for a given id
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
