package cli

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const agentNotesKey = "agent_notes"

// agentInheritedFlagNames is the persistent-flag allowlist shared by JSON help.
// Human --help still prints every flag. --yes is local to delete commands.
var agentInheritedFlagNames = map[string]bool{
	"account": true,
	"agent":   true,
	"json":    true,
	"md":      true,
	"jq":      true,
	"page":    true,
	"size":    true,
}

type agentHelpInfo struct {
	Command        string            `json:"command"`
	Path           string            `json:"path"`
	Short          string            `json:"short"`
	Long           string            `json:"long,omitempty"`
	Usage          string            `json:"usage"`
	Notes          []string          `json:"notes,omitempty"`
	Args           []agentArg        `json:"args,omitempty"`
	Subcommands    []agentSubcommand `json:"subcommands,omitempty"`
	Flags          []agentFlag       `json:"flags,omitempty"`
	InheritedFlags []agentFlag       `json:"inherited_flags,omitempty"`
}

type agentArg struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type agentSubcommand struct {
	Name  string `json:"name"`
	Short string `json:"short"`
	Path  string `json:"path"`
}

type agentFlag struct {
	Name      string   `json:"name"`
	Shorthand string   `json:"shorthand,omitempty"`
	Type      string   `json:"type"`
	Default   string   `json:"default"`
	Usage     string   `json:"usage"`
	Required  bool     `json:"required,omitempty"`
	Values    []string `json:"values,omitempty"`
	Format    string   `json:"format,omitempty"`
}

var useArgPattern = regexp.MustCompile(`(<[^>\s]+>|\[[^\]]+\])`)

func installAgentHelp(root *cobra.Command) {
	defaultHelp := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		agent, _ := cmd.Root().PersistentFlags().GetBool("agent")
		if agent {
			emitAgentHelp(cmd)
			return
		}
		defaultHelp(cmd, args)
	})
}

func emitAgentHelp(cmd *cobra.Command) {
	info := agentHelpInfo{
		Command:        cmd.Name(),
		Path:           cmd.CommandPath(),
		Short:          cmd.Short,
		Long:           strings.TrimSpace(cmd.Long),
		Usage:          cmd.UseLine(),
		Notes:          agentNotes(cmd),
		Args:           agentArgs(cmd.Use),
		Subcommands:    agentSubcommands(cmd),
		Flags:          agentFlags(cmd.NonInheritedFlags()),
		InheritedFlags: agentInheritedFlags(cmd),
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(info)
}

func agentNotes(cmd *cobra.Command) []string {
	raw := cmd.Annotations[agentNotesKey]
	if raw == "" {
		return nil
	}
	var notes []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			notes = append(notes, line)
		}
	}
	return notes
}

func agentArgs(use string) []agentArg {
	var args []agentArg
	for _, match := range useArgPattern.FindAllString(use, -1) {
		required := strings.HasPrefix(match, "<")
		name := strings.Trim(match, "<>[]")
		args = append(args, agentArg{Name: name, Required: required})
	}
	return args
}

func agentSubcommands(cmd *cobra.Command) []agentSubcommand {
	var subs []agentSubcommand
	parent := cmd.CommandPath()
	for _, child := range cmd.Commands() {
		if !child.IsAvailableCommand() {
			continue
		}
		subs = append(subs, agentSubcommand{
			Name:  child.Name(),
			Short: child.Short,
			Path:  child.CommandPath(),
		})
		for _, alias := range child.Aliases {
			subs = append(subs, agentSubcommand{
				Name:  alias,
				Short: child.Short,
				Path:  strings.TrimSpace(parent + " " + alias),
			})
		}
	}
	return subs
}

func agentFlags(set *pflag.FlagSet) []agentFlag {
	if set == nil {
		return nil
	}
	var flags []agentFlag
	set.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		flags = append(flags, flagInfo(f))
	})
	return flags
}

func agentInheritedFlags(cmd *cobra.Command) []agentFlag {
	var flags []agentFlag
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || !agentInheritedFlagNames[f.Name] {
			return
		}
		flags = append(flags, flagInfo(f))
	})
	return flags
}

func flagInfo(f *pflag.Flag) agentFlag {
	info := agentFlag{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Type:      f.Value.Type(),
		Default:   f.DefValue,
		Usage:     f.Usage,
	}
	if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
		info.Required = true
	}
	if values := f.Annotations["linden_enum"]; len(values) > 0 {
		info.Values = values
	}
	if formats := f.Annotations["linden_format"]; len(formats) > 0 {
		info.Format = formats[0]
	}
	return info
}
