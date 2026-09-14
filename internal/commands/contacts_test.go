package commands

import "testing"

func TestContactsCommandShape(t *testing.T) {
	t.Parallel()

	cmd := NewContactsCmd()
	want := map[string]bool{
		"list": false, "show": false, "search": false, "types": false,
	}
	for _, child := range cmd.Commands() {
		if _, ok := want[child.Name()]; ok {
			want[child.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing contacts %s command", name)
		}
	}

	list, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Fatal(err)
	}
	if list.Flags().Lookup("page") == nil || list.Flags().Lookup("size") == nil {
		t.Error("contacts list must define --page and --size")
	}

	search, _, err := cmd.Find([]string{"search"})
	if err != nil {
		t.Fatal(err)
	}
	q := search.Flags().Lookup("q")
	if q == nil {
		t.Fatal("contacts search must define --q")
	}
	if annotation := q.Annotations["cobra_annotation_bash_completion_one_required_flag"]; len(annotation) == 0 {
		t.Error("contacts search --q must be required")
	}
}
