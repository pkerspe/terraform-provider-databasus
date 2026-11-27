// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SettingsResponseModel struct {
	Id                                string `json:"id"`
	IsAllowExternalRegistrations      bool   `json:"isAllowExternalRegistrations"`
	IsAllowMemberInvitations          bool   `json:"isAllowMemberInvitations"`
	IsMemberAllowedToCreateWorkspaces bool   `json:"isMemberAllowedToCreateWorkspaces"`
	IsTwoFactorAuthRequired           bool   `json:"isTwoFactorAuthRequired"`
	IsEmailConfigured                 bool   `json:"isEmailConfigured"`
	TwoFactorRequirement              string `json:"twoFactorRequirement"`
}

type SettingsDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
	TwoFactorRequirement              types.String `tfsdk:"two_factor_requirement"`
}

type SettingsResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	IsAllowExternalRegistrations      types.Bool   `tfsdk:"allow_external_registrations"`
	IsAllowMemberInvitations          types.Bool   `tfsdk:"allow_member_invitations"`
	IsMemberAllowedToCreateWorkspaces types.Bool   `tfsdk:"member_allowed_to_create_workspaces"`
	TwoFactorRequirement              types.String `tfsdk:"two_factor_requirement"`
}

func (c *DatabasusClient) CreateUsersSettings(ctx context.Context, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool, twoFactorRequirement string) (*SettingsResponseModel, error) {
	return c.UpdateUsersSettings(ctx, "", allowExternalRegistrations, allowMemberInvitations, memberAllowedToCreateWorkspaces, twoFactorRequirement)
}

func (c *DatabasusClient) GetUsersSettings(ctx context.Context) (*SettingsResponseModel, *ErrorDetails) {
	var result SettingsResponseModel
	if err := c.doRequest(ctx, "GET", "/instance-settings", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) UpdateUsersSettings(ctx context.Context, _ string, allowExternalRegistrations bool, allowMemberInvitations bool, memberAllowedToCreateWorkspaces bool, twoFactorRequirement string) (*SettingsResponseModel, error) {
	current, err := c.GetUsersSettings(ctx)
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"isAllowExternalRegistrations":      allowExternalRegistrations,
		"isAllowMemberInvitations":          allowMemberInvitations,
		"isMemberAllowedToCreateWorkspaces": memberAllowedToCreateWorkspaces,
		"isTwoFactorAuthRequired":           current.IsTwoFactorAuthRequired,
		"twoFactorRequirement":              twoFactorRequirement,
	}
	b, _ := json.Marshal(body)
	var result SettingsResponseModel
	if err := c.doRequest(ctx, "PUT", "/instance-settings", bytes.NewBuffer(b), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *DatabasusClient) DeleteUsersSettings(ctx context.Context, id string) error { return nil }
