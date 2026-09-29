// Package skills_test is a contract suite for Agent Skills under skills/.
//
// It checks that the linden hub and its references exist, that SKILL.md has
// valid frontmatter, and that every documented `linden …` command is a real
// CLI path. Run with `make test-skills` or `go test ./skills`.
// These tests stay in the monorepo; publish-skills.sh strips *_test.go from
// the public skills package.
package skills_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mylinden-tech/linden-cli/internal/cli"
)

var requiredRefs = []string{
	"linden/references/envelope.md",
	"linden/references/auth-and-accounts.md",
	"linden/references/doctor.md",
	"linden/references/persons.md",
	"linden/references/persons-examples.md",
	"linden/references/contacts.md",
	"linden/references/memberships.md",
	"linden/references/invitations.md",
	"linden/references/account-shares.md",
	"linden/references/account-settings.md",
	"linden/references/pets.md",
	"linden/references/reminders.md",
	"linden/references/todos.md",
	"linden/references/vehicles.md",
	"linden/references/real-estates.md",
	"linden/references/online-accounts.md",
	"linden/references/insurances.md",
	"linden/references/wills.md",
	"linden/references/share-links.md",
}

var forbidden = []string{
	"linden properties",
	"linden documents",
	"linden setup",
	"curl ",
}

var cmdRe = regexp.MustCompile("(?m)(?:^|[ \t]|`)linden(?:[ \t]+\\\\)?[ \t]+([a-z][a-z0-9-]*)(?:[ \t]+([a-z][a-z0-9-]*))?")
var refRe = regexp.MustCompile(`references/[a-z0-9-]+\.md`)

func TestOnlyLindenSkill(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "linden") {
			continue
		}
		if e.Name() != "linden" {
			t.Errorf("unexpected skill directory %s", e.Name())
		}
	}
	if _, err := os.Stat(filepath.Join("linden", "SKILL.md")); err != nil {
		t.Errorf("missing linden/SKILL.md: %v", err)
	}
	for _, rel := range requiredRefs {
		if _, err := os.Stat(rel); err != nil {
			t.Errorf("missing %s", rel)
		}
	}
}

func TestFrontmatter(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("linden", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	name, desc, ok := parseFrontmatter(body)
	if !ok {
		t.Fatal("missing YAML frontmatter with name and description")
	}
	if name != "linden" {
		t.Errorf("name %q", name)
	}
	if !strings.Contains(desc, "Use when") {
		t.Errorf("description must start from %q", "Use when")
	}
	if strings.Contains(body, "argument-hint:") && strings.Contains(frontmatterBlock(body), "<") {
		t.Error("argument-hint and description must not use angle brackets")
	}
	if !strings.Contains(frontmatterBlock(body), `argument-hint: "[domain]"`) {
		t.Error(`frontmatter must set argument-hint: "[domain]"`)
	}
	if c := strings.Count(body, "\n"); c > 120 {
		t.Errorf("SKILL.md has %d lines; keep at 120 or fewer", c)
	}
	for _, f := range forbidden {
		if strings.Contains(strings.ToLower(body), f) {
			t.Errorf("forbidden string %q", f)
		}
	}
	for _, rel := range refRe.FindAllString(body, -1) {
		path := filepath.Join("linden", rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("SKILL.md mentions missing %s", path)
		}
	}
}

func TestDocumentedCommandsExist(t *testing.T) {
	seen := map[string]string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Error(readErr)
			return nil
		}
		body := stripFrontmatter(string(b))
		for _, f := range forbidden {
			if strings.Contains(strings.ToLower(body), f) {
				t.Errorf("%s: forbidden string %q", path, f)
			}
		}
		for _, m := range cmdRe.FindAllStringSubmatch(body, -1) {
			if strings.Contains(m[0], "<") || strings.Contains(m[0], "…") || strings.Contains(m[0], "...") {
				continue
			}
			key := m[1]
			if m[2] != "" && !strings.HasPrefix(m[2], "--") {
				key = m[1] + " " + m[2]
			}
			if _, ok := seen[key]; !ok {
				seen[key] = path
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(seen) == 0 {
		t.Fatal("no linden commands found in skills")
	}
	for key, path := range seen {
		args := strings.Fields(key)
		args = append(args, "--help")
		if err := runHelp(args); err != nil {
			t.Errorf("%s: %q is not a CLI command: %v", path, key, err)
		}
	}
}

func runHelp(args []string) error {
	root := cli.NewRootCmd(cli.BuildInfo{Version: "test"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	return root.Execute()
}

func stripFrontmatter(body string) string {
	if !strings.HasPrefix(body, "---\n") {
		return body
	}
	rest := strings.TrimPrefix(body, "---\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return body
	}
	return rest[end+4:]
}

func frontmatterBlock(body string) string {
	if !strings.HasPrefix(body, "---\n") {
		return ""
	}
	rest := strings.TrimPrefix(body, "---\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func parseFrontmatter(body string) (name, desc string, ok bool) {
	fm := frontmatterBlock(body)
	if fm == "" {
		return "", "", false
	}
	lines := strings.Split(fm, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "name:") {
			name = strings.TrimSpace(strings.Trim(strings.TrimPrefix(line, "name:"), `"'`))
			continue
		}
		if !strings.HasPrefix(line, "description:") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		if raw == "|" || raw == ">" || raw == "" {
			var b strings.Builder
			for _, extra := range lines[i+1:] {
				if extra == "" {
					b.WriteString(" ")
					continue
				}
				if extra[0] != ' ' && extra[0] != '\t' {
					break
				}
				if b.Len() > 0 {
					b.WriteString(" ")
				}
				b.WriteString(strings.TrimSpace(extra))
			}
			desc = strings.TrimSpace(b.String())
		} else {
			desc = strings.Trim(raw, `"'`)
		}
	}
	return name, desc, name != "" && desc != ""
}
