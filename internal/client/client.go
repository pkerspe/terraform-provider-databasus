// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Debug   bool
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

func NewDatabasusClient(baseURL, token string, verifySsl bool) *DatabasusClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !verifySsl,
		},
	}
	// Create a new HTTP client with the custom RoundTripper
	client := &http.Client{
		Transport: &CustomRoundTripper{
			Transport: transport,
			Token:     token,
		},
	}

	debug := os.Getenv("TF_DATABASUS_DEBUG") == "1"

	return &DatabasusClient{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    client,
		Debug:   debug,
	}
}

func GetJWT(baseURL, email, password string, verifySsl bool) (string, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	b, _ := json.Marshal(body)

	// Create a new HTTP client with SSL verification disabled if verifySsl is false
	httpClient := &http.Client{
		Transport: &CustomRoundTripper{
			Transport: http.DefaultTransport,
			Token:     "",
		},
	}

	if !verifySsl {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

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

	// Buffer the request body so it can be both logged and sent.
	var reqBodyBytes []byte
	if body != nil {
		var err error
		reqBodyBytes, err = io.ReadAll(body)
		if err != nil {
			errorDetails.ErrorInst = fmt.Errorf("failed to read request body: %w", err)
			return &errorDetails
		}
	}

	if c.Debug {
		fmt.Fprintf(os.Stderr, "[DATABASUS DEBUG] --> %s %s\n", method, url)
		if len(reqBodyBytes) > 0 {
			fmt.Fprintf(os.Stderr, "[DATABASUS DEBUG]     Request body: %s\n", string(reqBodyBytes))
		}
	}

	var reqReader io.Reader
	if len(reqBodyBytes) > 0 {
		reqReader = bytes.NewReader(reqBodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqReader)
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

	// Read the full response body so it can be logged and/or decoded.
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errorDetails.ErrorInst = fmt.Errorf("failed to read response body: %w", err)
		return &errorDetails
	}

	if c.Debug {
		fmt.Fprintf(os.Stderr, "[DATABASUS DEBUG] <-- %d %s\n", resp.StatusCode, url)
		if len(respBodyBytes) > 0 {
			fmt.Fprintf(os.Stderr, "[DATABASUS DEBUG]     Response body: %s\n", string(respBodyBytes))
		}
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorDetails.ResponseBody = string(respBodyBytes)
		errorDetails.ErrorCode = resp.StatusCode
		errorDetails.ErrorInst = fmt.Errorf("API error: status=%d body=%s", resp.StatusCode, string(respBodyBytes))
		return &errorDetails
	}

	// Decode if output is provided
	if out != nil {
		if err := json.NewDecoder(bytes.NewReader(respBodyBytes)).Decode(out); err != nil {
			errorDetails.ErrorInst = fmt.Errorf("failed to decode response: %w", err)
			return &errorDetails
		}
	}

	return nil
}
