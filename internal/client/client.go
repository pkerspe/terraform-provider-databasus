// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type ErrorDetails struct {
	ErrorCode    int
	ResponseBody string
	ErrorInst    error
}

func (e *ErrorDetails) Error() string {
	return e.ErrorInst.Error()
}

func (e *ErrorDetails) IsNotFound() bool {
	return e.ErrorCode == 404 ||
		strings.Contains(strings.ToLower(e.ResponseBody), "record not found") //fix for invalid API RCs from Databasus
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

func (c *DatabasusClient) doRequest(ctx context.Context, method, path string, body io.Reader, out interface{}) *ErrorDetails {
	var errorDetails ErrorDetails
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		errorDetails.ErrorInst = fmt.Errorf("failed to create request: %w", err)
		return &errorDetails
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		errorDetails.ErrorInst = fmt.Errorf("request failed: %w", err)
		return &errorDetails
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		errorDetails.ResponseBody = string(respBody)
		errorDetails.ErrorCode = resp.StatusCode
		errorDetails.ErrorInst = fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return &errorDetails
		//return fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBody)), &errorDetails
	}

	// Decode if output is provided
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			errorDetails.ErrorInst = fmt.Errorf("failed to decode response: %w", err)
			return &errorDetails
		}
	}

	return nil
}
