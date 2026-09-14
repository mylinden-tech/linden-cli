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

func TestRedactForAgentFailClosedOnMarshalError(t *testing.T) {
	t.Parallel()

	ch := make(chan int)
	out := output.RedactForAgent(ch, "secret")
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", out)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map on marshal error, got %v", m)
	}
}

func TestRedactForAgentFailClosedOnFuncMarshalError(t *testing.T) {
	t.Parallel()

	out := output.RedactForAgent(func() {}, "secret")
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", out)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map on marshal error, got %v", m)
	}
}

func TestRedactForAgentNoFieldsReturnsOriginal(t *testing.T) {
	t.Parallel()

	in := map[string]any{"vin": "SECRET"}
	out := output.RedactForAgent(in)
	if out == nil {
		t.Fatal("expected original data when no fields")
	}
	m, ok := out.(map[string]any)
	if !ok || m["vin"] != "SECRET" {
		t.Fatalf("expected unchanged data when no fields, got %v", out)
	}
}

func TestRedactForAgentNilDataReturnsNil(t *testing.T) {
	t.Parallel()

	if out := output.RedactForAgent(nil, "vin"); out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}
