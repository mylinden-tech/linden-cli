package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/client"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewPetsCmd creates the pets command group.
func NewPetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pets",
		Aliases: []string{"pet"},
		Short:   "Manage pets",
		Long:    "List, view, create, update, and delete pets in the active account.",
	}
	cmd.AddCommand(
		newPetsListCmd(),
		newPetsShowCmd(),
		newPetsCreateCmd(),
		newPetsUpdateCmd(),
		newPetsDeleteCmd(),
	)
	return cmd
}

func newPetsListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pets in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			if err := app.RequireAccount(); err != nil {
				return err
			}

			pets, err := app.Client.ListPets(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}

			count := len(pets)
			noun := "pets"
			if count == 1 {
				noun = "pet"
			}

			return app.OK(pets,
				output.WithSummary(fmt.Sprintf("%d %s", count, noun)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden pets show <id>",
						Description: "View pet details",
					},
					output.Breadcrumb{
						Action:      "create",
						Cmd:         "linden pets create",
						Description: "Create a new pet",
					},
				),
			)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newPetsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a pet by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			pet, err := app.Client.GetPet(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			return app.OK(pet,
				output.WithSummary(pet.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "update",
						Cmd:         "linden pets update " + pet.ID,
						Description: "Update this pet",
					},
				),
			)
		},
	}
}

func newPetsCreateCmd() *cobra.Command {
	var (
		name, species, breed, color, notes string
		birthDate, microchip, weight       string
		isActive                           bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new pet",
		Long: `Create a new pet in the active account.

  linden pets create --name "Rex" --species dog --breed "Golden Retriever"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			if err := app.RequireAccount(); err != nil {
				return err
			}

			req := client.PetCreateRequest{
				Name:            name,
				Species:         species,
				Breed:           strPtr(breed),
				Color:           strPtr(color),
				Notes:           strPtr(notes),
				BirthDate:       strPtr(birthDate),
				MicrochipNumber: strPtr(microchip),
				Weight:          strPtr(weight),
				IsActive:        isActive,
			}

			pet, err := app.Client.CreatePet(
				cmd.Context(), app.Config.AccountID, req,
			)
			if err != nil {
				return err
			}

			return app.OK(pet,
				output.WithSummary("Created "+pet.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden pets show " + pet.ID,
						Description: "View created pet",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Pet name (required)")
	cmd.Flags().StringVar(&species, "species", "", "Species (required)")
	cmd.Flags().StringVar(&breed, "breed", "", "Breed")
	cmd.Flags().StringVar(&color, "color", "", "Color")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&birthDate, "birth-date", "", "Birth date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&microchip, "microchip", "", "Microchip number")
	cmd.Flags().StringVar(&weight, "weight", "", "Weight")
	cmd.Flags().BoolVar(&isActive, "is-active", true, "Mark pet as active")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("species")
	return cmd
}

func newPetsUpdateCmd() *cobra.Command {
	var (
		name, species, breed, color, notes string
		birthDate, microchip, weight       string
		isActive                           bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a pet",
		Long: `Update one or more fields on a pet. Only flags you pass are sent to the API.

  linden pets update <id> --name "Rocket" --weight "25 kg"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			req := client.PetUpdateRequest{}
			changed := false

			if cmd.Flags().Changed("name") {
				req.Name = &name
				changed = true
			}
			if cmd.Flags().Changed("species") {
				req.Species = &species
				changed = true
			}
			if cmd.Flags().Changed("breed") {
				req.Breed = strPtr(breed)
				changed = true
			}
			if cmd.Flags().Changed("color") {
				req.Color = strPtr(color)
				changed = true
			}
			if cmd.Flags().Changed("notes") {
				req.Notes = strPtr(notes)
				changed = true
			}
			if cmd.Flags().Changed("birth-date") {
				req.BirthDate = strPtr(birthDate)
				changed = true
			}
			if cmd.Flags().Changed("microchip") {
				req.MicrochipNumber = strPtr(microchip)
				changed = true
			}
			if cmd.Flags().Changed("weight") {
				req.Weight = strPtr(weight)
				changed = true
			}
			if cmd.Flags().Changed("is-active") {
				req.IsActive = boolPtr(isActive)
				changed = true
			}

			if !changed {
				return noChanges(cmd)
			}

			pet, err := app.Client.UpdatePet(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}

			return app.OK(pet,
				output.WithSummary("Updated "+pet.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden pets show " + pet.ID,
						Description: "View updated pet",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Pet name")
	cmd.Flags().StringVar(&species, "species", "", "Species")
	cmd.Flags().StringVar(&breed, "breed", "", "Breed")
	cmd.Flags().StringVar(&color, "color", "", "Color")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&birthDate, "birth-date", "", "Birth date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&microchip, "microchip", "", "Microchip number")
	cmd.Flags().StringVar(&weight, "weight", "", "Weight")
	cmd.Flags().BoolVar(&isActive, "is-active", true, "Set whether pet is active")
	return cmd
}

func newPetsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a pet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			petID := args[0]

			if !yes {
				return output.ErrUsage("pass --yes to confirm deletion")
			}

			if err := app.Client.DeletePet(cmd.Context(), petID); err != nil {
				return err
			}

			return app.OK(map[string]any{"id": petID, "deleted": true},
				output.WithSummary("Pet deleted"),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden pets list",
						Description: "List remaining pets",
					},
				),
			)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion")
	return cmd
}
