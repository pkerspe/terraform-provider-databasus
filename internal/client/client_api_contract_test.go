// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMarshallBackupConfigResourceModelUsesCurrentAPIShape(t *testing.T) {
	model := BackupConfigResourceModel{
		Enabled:                          types.BoolValue(true),
		DatabaseId:                       types.StringValue("database-id"),
		StorageId:                        types.StringValue("storage-id"),
		Interval:                         types.StringValue("WEEKLY"),
		TimeOfDay:                        types.StringValue("08:00"),
		Weekday:                          types.Int32Value(3),
		DayOfMonth:                       types.Int32Value(1),
		CronExpression:                   types.StringValue("0 0 * * *"),
		MaxFailedRetryCount:              types.Int32Value(2),
		Encryption:                       types.BoolValue(true),
		RetentionPolicyType:              types.StringValue("COUNT"),
		RetentionCount:                   types.Int32Value(30),
		RetentionTimePeriod:              types.StringValue("MONTH"),
		RetentionGfsHours:                types.Int32Value(24),
		RetentionGfsDays:                 types.Int32Value(14),
		RetentionGfsWeeks:                types.Int32Value(4),
		RetentionGfsMonths:               types.Int32Value(12),
		RetentionGfsYears:                types.Int32Value(3),
		SendNotificationsOnBackupSuccess: types.BoolValue(false),
		SendNotificationsOnBackupFailure: types.BoolValue(true),
	}

	encoded, err := json.Marshal(marshallBackupConfigResourceModel(model))
	if err != nil {
		t.Fatalf("marshal backup config: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(encoded, &body); err != nil {
		t.Fatalf("unmarshal backup config: %v", err)
	}
	if body["storageId"] != "storage-id" || body["encryption"] != "ENCRYPTED" {
		t.Fatalf("unexpected common config fields: %#v", body)
	}
	settings, ok := body["settings"].(map[string]any)
	if !ok || settings["type"] != "LOGICAL" || settings["maxFailedTriesCount"] != float64(2) {
		t.Fatalf("unexpected settings: %#v", body["settings"])
	}
	if _, ok := body["databaseId"]; ok {
		t.Fatal("databaseId belongs in the URL, not the request body")
	}
}

func TestMapResponseToStorageS3ResourceModelSetsComputedStorageID(t *testing.T) {
	data := StorageS3ResourceModel{}
	MapResponseToStorageS3ResourceModel(&StorageS3ResponseModel{Id: "storage-id"}, &data)
	if got := data.StorageId.ValueString(); got != "storage-id" {
		t.Fatalf("storage_id = %q, want %q", got, "storage-id")
	}
}
