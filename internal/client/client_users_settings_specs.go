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
*			SETTINGS CRUD functions and model definitions
********************************************************/

type SettingsResponseModel struct {
	Id                                string `json:"id"`
	IsAllowExternalRegistrations      bool   `json:"isAllowExternalRegistrations"`
	IsAllowMemberInvitations          bool   `json:"isAllowMemberInvitations"`
	IsMemberAllowedToCreateWorkspaces bool   `json:"isMemberAllowedToCreateWorkspaces"`
}

type SettingsDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
}

type SettingsResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
}

// settings always exist, so the create just internally calls the update function.
func (c *DatabasusClient) CreateUsersSettings(ctx context.Context, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool) (*SettingsResponseModel, error) {
	currentSettings, err := c.GetUsersSettings(ctx)
	if err != nil {
		return nil, err
	}

	return c.UpdateUsersSettings(ctx, currentSettings.Id, allowExternalRegistrations, allowMemberInvitations, memberAllowedToCreateWorkspaces)
}

func (c *DatabasusClient) GetUsersSettings(ctx context.Context) (*SettingsResponseModel, *ErrorDetails) {
	var result SettingsResponseModel

	err := c.doRequest(ctx, "GET", "/users/settings", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DatabasusClient) UpdateUsersSettings(ctx context.Context, id string, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool) (*SettingsResponseModel, error) {
	var result SettingsResponseModel

	body := map[string]any{
		"isAllowExternalRegistrations":      allowExternalRegistrations,
		"isAllowMemberInvitations":          allowMemberInvitations,
		"isMemberAllowedToCreateWorkspaces": memberAllowedToCreateWorkspaces,
	}

	if id != "" {
		body["id"] = id
	}

	b, _ := json.Marshal(body)
	err := c.doRequest(ctx, "PUT", "/users/settings", bytes.NewBuffer(b), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete does not exist, so we just do nothing here.
func (c *DatabasusClient) DeleteUsersSettings(ctx context.Context, id string) error {
	return nil
}
