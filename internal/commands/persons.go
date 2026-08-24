package commands

import (
	"encoding/json"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/linden-family/linden-cli/internal/appctx"
	"github.com/linden-family/linden-cli/internal/client"
	"github.com/linden-family/linden-cli/internal/output"
	"github.com/linden-family/linden-cli/internal/tui"
	personsTUI "github.com/linden-family/linden-cli/internal/tui/persons"
)

// NewPersonsCmd creates the persons command group.
func NewPersonsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "persons",
		Aliases: []string{"person"},
		Short:   "Manage persons",
		Long:    "List, view, create, update, and delete persons in the active account.",
	}
	cmd.AddCommand(
		newPersonsListCmd(),
		newPersonsShowCmd(),
		newPersonsCreateCmd(),
		newPersonsUpdateCmd(),
		newPersonsDeleteCmd(),
	)
	return cmd
}

func newPersonsListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List persons in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			if err := app.RequireAccount(); err != nil {
				return err
			}

			persons, err := app.Client.ListPersons(cmd.Context(), app.Config.AccountID,
				client.ListPersonsParams{Page: page, Size: size})
			if err != nil {
				return err
			}

			// Launch TUI when interactive and not suppressed
			if app.IsInteractive() && os.Getenv("LINDEN_NO_TUI") == "" {
				return runPersonsTUI(app, persons)
			}

			count := len(persons)
			noun := "persons"
			if count == 1 {
				noun = "person"
			}

			return app.OK(persons,
				output.WithSummary(fmt.Sprintf("%d %s", count, noun)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden persons show <id>",
						Description: "View person details",
					},
					output.Breadcrumb{
						Action:      "create",
						Cmd:         "linden persons create",
						Description: "Create a new person",
					},
				),
			)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newPersonsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a person by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			person, err := app.Client.GetPerson(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			return app.OK(person,
				output.WithSummary(person.FirstName+" "+person.LastName),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "update",
						Cmd:         "linden persons update " + person.ID,
						Description: "Update this person",
					},
					output.Breadcrumb{
						Action:      "delete",
						Cmd:         "linden persons delete " + person.ID,
						Description: "Delete this person",
					},
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden persons list",
						Description: "Back to list",
					},
				),
			)
		},
	}
}

func newPersonsCreateCmd() *cobra.Command {
	var (
		firstName, lastName, email, phone, notes string
		birthday                                  string
		addressStreet, addressApt, addressCity    string
		addressState, addressZip, addressCountry  string
		relationshipType                          string
		isActive                                  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new person",
		Long: `Create a new person in the active account.

When run interactively without required flags, a form is shown.

  linden persons create --first-name "Jane" --last-name "Doe" --email jane@example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			if err := app.RequireAccount(); err != nil {
				return err
			}

			// Collect required fields interactively if not provided and in a TTY
			if firstName == "" || lastName == "" {
				if !app.IsInteractive() || os.Getenv("LINDEN_NO_TUI") != "" {
					return output.ErrUsage("--first-name and --last-name are required")
				}

				values, err := tui.Form("Create Person", []tui.FormField{
					{Key: "first_name", Title: "First name", Placeholder: "Jane", Required: true, Default: firstName},
					{Key: "last_name", Title: "Last name", Placeholder: "Doe", Required: true, Default: lastName},
					{Key: "email", Title: "Email", Placeholder: "jane@example.com", Default: email},
					{Key: "phone", Title: "Phone", Placeholder: "+1 555 000 0000", Default: phone},
					{Key: "relationship_type", Title: "Relationship", Placeholder: "spouse, sibling, parent…", Default: relationshipType},
				})
				if err != nil {
					return err
				}

				firstName = values["first_name"]
				lastName = values["last_name"]
				email = values["email"]
				phone = values["phone"]
				relationshipType = values["relationship_type"]
			}

			req := client.PersonCreateRequest{
				FirstName:        firstName,
				LastName:         lastName,
				Email:            strPtr(email),
				Phone:            strPtr(phone),
				Notes:            strPtr(notes),
				Birthday:         strPtr(birthday),
				IsActive:         isActive,
				AddressStreet:    strPtr(addressStreet),
				AddressApt:       strPtr(addressApt),
				AddressCity:      strPtr(addressCity),
				AddressState:     strPtr(addressState),
				AddressZip:       strPtr(addressZip),
				AddressCountry:   strPtr(addressCountry),
				RelationshipType: strPtr(relationshipType),
			}

			person, err := app.Client.CreatePerson(cmd.Context(), app.Config.AccountID, req)
			if err != nil {
				return err
			}

			return app.OK(person,
				output.WithSummary(fmt.Sprintf("Created %s %s", person.FirstName, person.LastName)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden persons show " + person.ID,
						Description: "View created person",
					},
					output.Breadcrumb{
						Action:      "update",
						Cmd:         "linden persons update " + person.ID,
						Description: "Update this person",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&firstName, "first-name", "", "First name (required)")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name (required)")
	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&birthday, "birthday", "", "Birthday (YYYY-MM-DD)")
	cmd.Flags().StringVar(&addressStreet, "address-street", "", "Street address")
	cmd.Flags().StringVar(&addressApt, "address-apt", "", "Apt/suite")
	cmd.Flags().StringVar(&addressCity, "address-city", "", "City")
	cmd.Flags().StringVar(&addressState, "address-state", "", "State/province")
	cmd.Flags().StringVar(&addressZip, "address-zip", "", "ZIP/postal code")
	cmd.Flags().StringVar(&addressCountry, "address-country", "", "Country")
	cmd.Flags().StringVar(&relationshipType, "relationship-type", "", "Relationship type")
	cmd.Flags().BoolVar(&isActive, "active", true, "Mark person as active")
	return cmd
}

func newPersonsUpdateCmd() *cobra.Command {
	var (
		firstName, lastName, email, phone, notes string
		birthday                                  string
		addressStreet, addressApt, addressCity    string
		addressState, addressZip, addressCountry  string
		relationshipType                          string
		setActive, setInactive                    bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a person",
		Long: `Update one or more fields on a person. Only flags you pass are sent to the API.

  linden persons update <id> --email new@example.com --phone "+1 555 111 2222"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			// Build update request from only the flags that were explicitly set
			req := client.PersonUpdateRequest{}
			changed := false

			if cmd.Flags().Changed("first-name") {
				req.FirstName = &firstName
				changed = true
			}
			if cmd.Flags().Changed("last-name") {
				req.LastName = &lastName
				changed = true
			}
			if cmd.Flags().Changed("email") {
				req.Email = strPtr(email)
				changed = true
			}
			if cmd.Flags().Changed("phone") {
				req.Phone = strPtr(phone)
				changed = true
			}
			if cmd.Flags().Changed("notes") {
				req.Notes = strPtr(notes)
				changed = true
			}
			if cmd.Flags().Changed("birthday") {
				req.Birthday = strPtr(birthday)
				changed = true
			}
			if cmd.Flags().Changed("address-street") {
				req.AddressStreet = strPtr(addressStreet)
				changed = true
			}
			if cmd.Flags().Changed("address-apt") {
				req.AddressApt = strPtr(addressApt)
				changed = true
			}
			if cmd.Flags().Changed("address-city") {
				req.AddressCity = strPtr(addressCity)
				changed = true
			}
			if cmd.Flags().Changed("address-state") {
				req.AddressState = strPtr(addressState)
				changed = true
			}
			if cmd.Flags().Changed("address-zip") {
				req.AddressZip = strPtr(addressZip)
				changed = true
			}
			if cmd.Flags().Changed("address-country") {
				req.AddressCountry = strPtr(addressCountry)
				changed = true
			}
			if cmd.Flags().Changed("relationship-type") {
				req.RelationshipType = strPtr(relationshipType)
				changed = true
			}
			if setActive {
				req.IsActive = boolPtr(true)
				changed = true
			}
			if setInactive {
				req.IsActive = boolPtr(false)
				changed = true
			}

			if !changed {
				return noChanges(cmd)
			}

			person, err := app.Client.UpdatePerson(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}

			return app.OK(person,
				output.WithSummary(fmt.Sprintf("Updated %s %s", person.FirstName, person.LastName)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden persons show " + person.ID,
						Description: "View updated person",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&birthday, "birthday", "", "Birthday (YYYY-MM-DD)")
	cmd.Flags().StringVar(&addressStreet, "address-street", "", "Street address")
	cmd.Flags().StringVar(&addressApt, "address-apt", "", "Apt/suite")
	cmd.Flags().StringVar(&addressCity, "address-city", "", "City")
	cmd.Flags().StringVar(&addressState, "address-state", "", "State/province")
	cmd.Flags().StringVar(&addressZip, "address-zip", "", "ZIP/postal code")
	cmd.Flags().StringVar(&addressCountry, "address-country", "", "Country")
	cmd.Flags().StringVar(&relationshipType, "relationship-type", "", "Relationship type")
	cmd.Flags().BoolVar(&setActive, "active", false, "Mark person as active")
	cmd.Flags().BoolVar(&setInactive, "inactive", false, "Mark person as inactive")
	return cmd
}

func newPersonsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a person",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			personID := args[0]

			if !yes && app.IsInteractive() && os.Getenv("LINDEN_NO_TUI") == "" {
				confirmed, err := tui.ConfirmDangerous(
					fmt.Sprintf("Delete person %q?", personID))
				if err != nil || !confirmed {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			} else if !yes {
				return output.ErrUsage("pass --yes to confirm deletion without a prompt")
			}

			if err := app.Client.DeletePerson(cmd.Context(), personID); err != nil {
				return err
			}

			return app.OK(map[string]any{"id": personID, "deleted": true},
				output.WithSummary("Person deleted"),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden persons list",
						Description: "List remaining persons",
					},
				),
			)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	return cmd
}

// runPersonsTUI launches the interactive bubbletea persons browser.
func runPersonsTUI(app *appctx.App, persons []client.Person) error {
	m := personsTUI.New()
	prog := tea.NewProgram(m)

	// Send loaded persons immediately from a goroutine
	go func() {
		prog.Send(personsTUI.LoadedMsg{Persons: persons})
	}()

	result, err := prog.Run()
	if err != nil {
		return err
	}

	final := result.(personsTUI.Model) //nolint:errcheck // type assertion always succeeds
	if sel := final.Selected(); sel != nil {
		b, err := json.MarshalIndent(sel, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}
