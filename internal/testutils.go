package internal

import (
	"encoding/json"
	"testing"
)

// TMustJSON is a helper to serialize data in json, fails the test if
// the serialization err
func TMustJSON(t *testing.T, data any) []byte {
	out, err := json.Marshal(data)
	if err != nil {
		t.Errorf("test setup failed: json serialization failed: %v", err)
		t.FailNow()
		return nil
	}

	return out
}
