package client

// MeasurementInput is a measured value as sent to the API: the value as
// entered (a decimal string, never a float) and its unit key, e.g. "kg".
type MeasurementInput struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// Measurement is a measured value as the API returns it: the value and unit as
// entered, plus the value in the viewer's unit system.
type Measurement struct {
	Value        string `json:"value"`
	Unit         string `json:"unit"`
	DisplayValue string `json:"display_value"`
	DisplayUnit  string `json:"display_unit"`
}
