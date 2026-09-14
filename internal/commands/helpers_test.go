package commands

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func testApp(agent bool) (*appctx.App, *bytes.Buffer) {
	var buf bytes.Buffer
	w, err := output.New(output.Options{
		Format: output.FormatJSON,
		Writer: &buf,
	})
	if err != nil {
		panic(err)
	}
	return &appctx.App{
		Output: w,
		Flags:  appctx.GlobalFlags{Agent: agent},
	}, &buf
}

func TestOkMaybeRedactedAgentMode(t *testing.T) {
	t.Parallel()

	app, buf := testApp(true)
	data := map[string]any{"id": "1", "vin": "SECRET", "make": "Toyota"}

	if err := okMaybeRedacted(app, data, output.VehicleAgentOmit); err != nil {
		t.Fatal(err)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	m, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", resp.Data)
	}
	if _, ok := m["vin"]; ok {
		t.Fatalf("vin should be redacted: %v", m)
	}
	if m["make"] != "Toyota" {
		t.Fatalf("make should remain: %v", m)
	}
}

func TestOkMaybeRedactedNonAgentMode(t *testing.T) {
	t.Parallel()

	app, buf := testApp(false)
	data := map[string]any{"id": "1", "vin": "SECRET"}

	if err := okMaybeRedacted(app, data, output.VehicleAgentOmit); err != nil {
		t.Fatal(err)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	m, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", resp.Data)
	}
	if m["vin"] != "SECRET" {
		t.Fatalf("vin should not be redacted without agent mode: %v", m)
	}
}

func TestOkMaybeRedactedEmptyFieldsSkipsRedaction(t *testing.T) {
	t.Parallel()

	app, buf := testApp(true)
	data := map[string]any{"vin": "SECRET"}

	if err := okMaybeRedacted(app, data, nil); err != nil {
		t.Fatal(err)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	m, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", resp.Data)
	}
	if m["vin"] != "SECRET" {
		t.Fatalf("empty fields should skip redaction: %v", m)
	}
}
