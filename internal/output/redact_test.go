package output_test

import (
	"encoding/json"
	"testing"

	"github.com/mylinden-tech/linden-cli/internal/output"
)

func TestRedactForAgentRemovesKeysFromMap(t *testing.T) {
	in := map[string]any{"id": "1", "vin": "SECRET", "make": "Toyota"}
	out := output.RedactForAgent(in, "vin")
	b, _ := json.Marshal(out)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["vin"]; ok {
		t.Fatalf("vin still present: %v", m)
	}
	if m["make"] != "Toyota" {
		t.Fatalf("make lost: %v", m)
	}
}

func TestRedactForAgentRemovesKeysFromSlice(t *testing.T) {
	in := []map[string]any{{"id": "1", "password": "x"}, {"id": "2", "password": "y"}}
	out := output.RedactForAgent(in, "password")
	b, _ := json.Marshal(out)
	var arr []map[string]any
	_ = json.Unmarshal(b, &arr)
	for i, m := range arr {
		if _, ok := m["password"]; ok {
			t.Fatalf("item %d still has password", i)
		}
	}
}
