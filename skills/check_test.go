// Package skills_test is a contract suite for Agent Skills under skills/.
//
// It checks that required skill trees and reference docs exist, that each
// SKILL.md has valid frontmatter, and that every documented `linden …`
// command is a real CLI path (and that known non-commands like documents or
// curl are absent). Run with `make test-skills` or `go test ./skills`.
// These tests stay in the monorepo; publish-skills.sh strips *_test.go from
// the public skills package.
package skills_test

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var requiredSkills = []string{
	"linden", "linden-account-settings", "linden-account-shares", "linden-contacts",
	"linden-doctor", "linden-invitations", "linden-memberships", "linden-persons",
	"linden-pets", "linden-reminders", "linden-todos",
	"linden-vehicles", "linden-real-estates", "linden-online-accounts",
	"linden-insurances", "linden-wills", "linden-share-links",
}

var hubSkills = map[string]bool{"linden": true, "linden-doctor": true}

var requiredRefs = []string{
	"linden/references/envelope.md",
	"linden/references/auth-and-accounts.md",
	"linden-persons/references/person-fields.md",
	"linden-persons/references/examples.md",
	"linden-pets/references/pet-fields.md",
}

// First two tokens after `linden` (command group + subcommand), or a single token for doctor.
var allowed = map[string]bool{
	"doctor":               true,
	"auth login":           true,
	"auth status":          true,
	"auth logout":          true,
	"accounts list":        true,
	"accounts show":        true,
	"accounts stats":       true,
	"accounts use":         true,
	"account-shares list":  true,
	"account-shares show":  true,
	"contacts list":        true,
	"contacts show":        true,
	"contacts search":      true,
	"contacts types":       true,
	"invitations list":     true,
	"invitations show":     true,
	"memberships list":     true,
	"memberships roles":    true,
	"memberships show":     true,
	"persons list":         true,
	"persons show":         true,
	"persons create":       true,
	"persons update":       true,
	"persons delete":       true,
	"pets list":            true,
	"pets show":            true,
	"pets create":          true,
	"pets update":          true,
	"pets delete":          true,
	"reminders list":       true,
	"reminders show":       true,
	"reminders create":     true,
	"reminders update":     true,
	"reminders complete":   true,
	"reminders delete":     true,
	"todos list":           true,
	"todos lists":          true,
	"todos create-list":    true,
	"todos show":           true,
	"todos create":         true,
	"todos update":         true,
	"todos delete":         true,
	"user-settings show":   true,
	"vehicles list":        true,
	"vehicles show":        true,
	"real-estates list":    true,
	"real-estates show":    true,
	"online-accounts list": true,
	"online-accounts show": true,
	"insurances list":      true,
	"insurances show":      true,
	"insurances providers": true,
	"insurances expiring":  true,
	"wills list":           true,
	"wills show":           true,
	"share-links list":     true,
	"share-links show":     true,
}

var forbidden = []string{
	"linden properties",
	"linden documents",
	"linden setup",
	"curl ",
}

var cmdRe = regexp.MustCompile("(?m)(?:^|[ \t]|`)linden(?:[ \t]+\\\\)?[ \t]+([a-z-]+)(?:[ \t]+([a-z-]+))?")

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

func TestDomainSkillsDisableModelInvocation(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "linden") {
			continue
		}
		if hubSkills[e.Name()] {
			continue
		}
		b, err := os.ReadFile(filepath.Join(e.Name(), "SKILL.md"))
		if err != nil {
			t.Errorf("%s: %v", e.Name(), err)
			continue
		}
		if !strings.Contains(string(b), "disable-model-invocation: true") {
			t.Errorf("%s: domain skill must set disable-model-invocation: true", e.Name())
		}
	}
}

func TestSkillLayout(t *testing.T) {
	root := "."
	for _, name := range requiredSkills {
		path := filepath.Join(root, name, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s", path)
		}
	}
	for _, rel := range requiredRefs {
		path := filepath.Join(root, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s", path)
		}
	}
}

func TestFrontmatterAndCommands(t *testing.T) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "SKILL.md" {
			return err
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Error(readErr)
			return nil
		}
		body := string(b)
		name, desc, ok := parseFrontmatter(body)
		if !ok {
			t.Errorf("%s: missing YAML frontmatter with name and description", path)
			return nil
		}
		if name == "" {
			t.Errorf("%s: empty name", path)
		}
		if !strings.Contains(desc, "This skill should be used when") {
			t.Errorf("%s: description must be third-person and include %q", path, "This skill should be used when")
		}
		if c := strings.Count(body, "\n"); c > 500 {
			t.Errorf("%s: SKILL.md has %d lines; keep under 500", path, c)
		}
		for _, f := range forbidden {
			if strings.Contains(strings.ToLower(body), f) {
				t.Errorf("%s: forbidden string %q", path, f)
			}
		}
		for _, m := range cmdRe.FindAllStringSubmatch(stripFrontmatter(body), -1) {
			group := m[1]
			if group == "--version" || group == "--help" || group == "--json" || group == "--agent" || group == "--jq" || group == "--md" {
				continue
			}
			key := group
			if m[2] != "" && !strings.HasPrefix(m[2], "--") {
				key = group + " " + m[2]
			}
			if group == "doctor" {
				key = "doctor"
			}
			if !allowed[key] && !allowed[group] {
				t.Errorf("%s: undocumented or invented command %q", path, strings.TrimSpace(m[0]))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReferenceLinksExist(t *testing.T) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		f, openErr := os.Open(path)
		if openErr != nil {
			t.Error(openErr)
			return nil
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		re := regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
		for sc.Scan() {
			for _, m := range re.FindAllStringSubmatch(sc.Text(), -1) {
				link := m[1]
				if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "#") {
					continue
				}
				target := filepath.Join(filepath.Dir(path), link)
				if _, statErr := os.Stat(target); statErr != nil {
					t.Errorf("%s: broken link %s", path, link)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func parseFrontmatter(body string) (name, desc string, ok bool) {
	if !strings.HasPrefix(body, "---\n") {
		return "", "", false
	}
	rest := strings.TrimPrefix(body, "---\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", "", false
	}
	fm := rest[:end]
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
		// Folded/literal block: description: | or description: >
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
