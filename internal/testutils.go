package internal

import (
	"encoding/json"
	"testing"
)

// MustJSON is a helper to serialize data in json, fails the test if
// the serialization err
func MustJSON(t *testing.T, data any) []byte {
	out, err := json.Marshal(data)
	if err != nil {
		t.Errorf("test setup failed: json serialization failed: %v", err)
		t.FailNow()
		return nil
	}

	return out
}
