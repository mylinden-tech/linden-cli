package commands

import (
	"regexp"
	"strings"

	"github.com/mylinden-tech/linden-cli/internal/client"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// massUnits are the weight units the API accepts.
var massUnits = map[string]bool{"kg": true, "g": true, "lb": true, "oz": true}

// weightPattern is a number (decimal point or comma) followed by a unit,
// with or without a space: "25 lb", "12,5kg".
var weightPattern = regexp.MustCompile(`^(\d+(?:[.,]\d+)?)\s*([A-Za-z]+)$`)

const weightHelp = `Weight with its unit: kg, g, lb or oz (e.g. "12.5 kg", "25 lb")`

// parseWeight reads a --weight flag such as "25 lb" into the API's
// {value, unit} shape. A weight needs its unit: the API never guesses one.
func parseWeight(text string) (*client.MeasurementInput, error) {
	m := weightPattern.FindStringSubmatch(strings.TrimSpace(text))
	if m == nil || !massUnits[strings.ToLower(m[2])] {
		return nil, output.ErrUsage(
			`--weight takes a number and a unit (kg, g, lb or oz), e.g. "25 lb"`,
		)
	}
	return &client.MeasurementInput{
		Value: strings.ReplaceAll(m[1], ",", "."),
		Unit:  strings.ToLower(m[2]),
	}, nil
}
