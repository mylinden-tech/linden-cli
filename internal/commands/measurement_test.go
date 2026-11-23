package commands

import (
	"testing"

	"github.com/mylinden-tech/linden-cli/internal/client"
)

func TestParseWeight(t *testing.T) {
	t.Parallel()

	accepted := map[string]client.MeasurementInput{
		"25 lb":   {Value: "25", Unit: "lb"},
		"25lb":    {Value: "25", Unit: "lb"},
		"12.5 kg": {Value: "12.5", Unit: "kg"},
		"12,5 KG": {Value: "12.5", Unit: "kg"},
		" 800 g ": {Value: "800", Unit: "g"},
		"3 oz":    {Value: "3", Unit: "oz"},
	}
	for text, want := range accepted {
		got, err := parseWeight(text)
		if err != nil {
			t.Errorf("parseWeight(%q) error = %v", text, err)
			continue
		}
		if *got != want {
			t.Errorf("parseWeight(%q) = %+v, want %+v", text, *got, want)
		}
	}

	// No unit is never guessed; other units and malformed values are rejected.
	for _, text := range []string{"30", "", "about 30 kg", "30 m", "-5 kg", "25 lbs"} {
		if got, err := parseWeight(text); err == nil {
			t.Errorf("parseWeight(%q) = %+v, want an error", text, got)
		}
	}
}
