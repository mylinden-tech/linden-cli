package commands

import "testing"

func TestAccountsCommandIncludesReadOnlyDetails(t *testing.T) {
	t.Parallel()

	cmd := NewAccountsCmd()
	for _, name := range []string{"show", "stats"} {
		child, _, err := cmd.Find([]string{name})
		if err != nil || child.Name() != name {
			t.Errorf("missing accounts %s command", name)
		}
	}

	show, _, err := cmd.Find([]string{"show"})
	if err != nil {
		t.Fatal(err)
	}
	if err := show.Args(show, nil); err != nil {
		t.Errorf("accounts show should accept no ID: %v", err)
	}
	if err := show.Args(show, []string{"account-id"}); err != nil {
		t.Errorf("accounts show should accept one ID: %v", err)
	}
}

func TestUserSettingsCommandShape(t *testing.T) {
	t.Parallel()

	cmd := NewUserSettingsCmd()
	show, _, err := cmd.Find([]string{"show"})
	if err != nil || show.Name() != "show" {
		t.Fatal("missing user-settings show command")
	}
}
