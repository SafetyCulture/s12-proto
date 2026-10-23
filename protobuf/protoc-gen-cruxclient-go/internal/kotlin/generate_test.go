// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import "testing"

func TestEnumEntryNameCollision(t *testing.T) {
	names := map[string]string{}
	if err := checkKotlinName(names, "name", entryName("name")); err != nil {
		t.Fatal(err)
	}
	err := checkKotlinName(names, "name_", entryName("name_"))
	if err == nil || err.Error() != `name and name_ map to the same Kotlin name "name_"` {
		t.Fatalf("error = %v", err)
	}
}
