package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestAgentHelpJSON(t *testing.T) {
	root := NewRootCmd(BuildInfo{Version: "test"})
	walkCommands(root, func(cmd *cobra.Command) {
		t.Run(cmd.CommandPath(), func(t *testing.T) {
			out := runHelp(t, append(commandArgs(cmd), "--agent", "--help"))
			var info agentHelpInfo
			if err := json.Unmarshal([]byte(out), &info); err != nil {
				t.Fatalf("invalid JSON: %v\n%s", err, out)
			}
			if info.Path != cmd.CommandPath() {
				t.Fatalf("path %q, want %q", info.Path, cmd.CommandPath())
			}
		})
	})
}

func TestJSONHelpStaysText(t *testing.T) {
	out := runHelp(t, []string{"persons", "list", "--json", "--help"})
	if strings.Contains(out, `"inherited_flags"`) {
		t.Fatalf("--json --help returned agent JSON:\n%s", out)
	}
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("expected human help, got:\n%s", out)
	}
}

func TestRequiredFlagsMarked(t *testing.T) {
	root := NewRootCmd(BuildInfo{Version: "test"})
	walkCommands(root, func(cmd *cobra.Command) {
		cmd.NonInheritedFlags().VisitAll(func(f *pflag.Flag) {
			if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
				return
			}
			info := agentHelp(t, cmd)
			for _, flag := range info.Flags {
				if flag.Name == f.Name {
					if !flag.Required {
						t.Errorf("%s --%s: required is false", cmd.CommandPath(), f.Name)
					}
					return
				}
			}
			t.Errorf("%s --%s: missing from JSON flags", cmd.CommandPath(), f.Name)
		})
	})
}

func TestTodoStatusValues(t *testing.T) {
	root := NewRootCmd(BuildInfo{Version: "test"})
	update, _, err := root.Find([]string{"todos", "update"})
	if err != nil {
		t.Fatal(err)
	}
	info := agentHelp(t, update)
	for _, flag := range info.Flags {
		if flag.Name != "status" {
			continue
		}
		if strings.Join(flag.Values, ",") != "open,completed" {
			t.Fatalf("status values %v", flag.Values)
		}
		return
	}
	t.Fatal("todos update is missing --status")
}

func TestDeleteYesIsLocal(t *testing.T) {
	root := NewRootCmd(BuildInfo{Version: "test"})
	del, _, err := root.Find([]string{"persons", "delete"})
	if err != nil {
		t.Fatal(err)
	}
	info := agentHelp(t, del)
	for _, flag := range info.InheritedFlags {
		if flag.Name == "yes" {
			t.Fatal("--yes must not be an inherited flag")
		}
	}
	for _, flag := range info.Flags {
		if flag.Name == "yes" {
			return
		}
	}
	t.Fatal("persons delete is missing local --yes")
}

func TestReminderCreateNotes(t *testing.T) {
	root := NewRootCmd(BuildInfo{Version: "test"})
	create, _, err := root.Find([]string{"reminders", "create"})
	if err != nil {
		t.Fatal(err)
	}
	info := agentHelp(t, create)
	if len(info.Notes) != 1 || !strings.Contains(info.Notes[0], "--person") || !strings.Contains(info.Notes[0], "--pet") {
		t.Fatalf("notes %v", info.Notes)
	}
}

func agentHelp(t *testing.T, cmd *cobra.Command) agentHelpInfo {
	t.Helper()
	out := runHelp(t, append(commandArgs(cmd), "--agent", "--help"))
	var info agentHelpInfo
	if err := json.Unmarshal([]byte(out), &info); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	return info
}

func runHelp(t *testing.T, args []string) string {
	t.Helper()
	root := NewRootCmd(BuildInfo{Version: "test"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(args, " "), err, buf.String())
	}
	return buf.String()
}

func commandArgs(cmd *cobra.Command) []string {
	path := strings.Fields(cmd.CommandPath())
	if len(path) > 0 && path[0] == "linden" {
		path = path[1:]
	}
	return path
}

func walkCommands(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, child := range cmd.Commands() {
		if child.Name() == "help" || !child.IsAvailableCommand() {
			continue
		}
		walkCommands(child, fn)
	}
}
