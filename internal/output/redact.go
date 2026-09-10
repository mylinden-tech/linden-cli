package output

import "encoding/json"

var (
	VehicleAgentOmit       = []string{"vin", "license_plate"}
	OnlineAccountAgentOmit = []string{"password"}
	InsuranceAgentOmit     = []string{"policy_number"}
	WillAgentOmit          = []string{"document_file_url", "document_file_id", "notes", "plain_language_summary"}
)

func RedactForAgent(data any, fields ...string) any {
	if len(fields) == 0 || data == nil {
		return data
	}
	b, err := json.Marshal(data)
	if err != nil {
		return map[string]any{}
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return map[string]any{}
	}
	omit := map[string]struct{}{}
	for _, f := range fields {
		omit[f] = struct{}{}
	}
	return stripKeys(v, omit)
}

func stripKeys(v any, omit map[string]struct{}) any {
	switch t := v.(type) {
	case map[string]any:
		for k := range omit {
			delete(t, k)
		}
		for k, child := range t {
			t[k] = stripKeys(child, omit)
		}
		return t
	case []any:
		for i, child := range t {
			t[i] = stripKeys(child, omit)
		}
		return t
	default:
		return v
	}
}
