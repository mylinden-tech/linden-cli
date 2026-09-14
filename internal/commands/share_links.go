package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewShareLinksCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "share-links", Aliases: []string{"share-link"}, Short: "View resource share links", Long: "List and view share links for a resource. This is distinct from account-shares."}
	cmd.AddCommand(newShareLinksListCmd(), newShareLinksShowCmd())
	return cmd
}

func newShareLinksListCmd() *cobra.Command {
	var resourceType, resourceID string
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List share links for a resource", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		id, err := resolveShareResourceID(app, resourceID)
		if err != nil {
			return err
		}
		result, err := app.Client.ListShareLinks(cmd.Context(), resourceType, id, page, size)
		if err != nil {
			return err
		}
		return app.OK(result.Items,
			output.WithSummary(countSummary(len(result.Items), result.Total, "share link", "share links")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: fmt.Sprintf("linden share-links show <id> --resource-type %s --resource-id %s", resourceType, id), Description: "View share link details"}))
	}}
	addShareResourceFlags(cmd, &resourceType, &resourceID)
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newShareLinksShowCmd() *cobra.Command {
	var resourceType, resourceID string
	cmd := &cobra.Command{Use: "show <id>", Short: "Show a share link by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		id, err := resolveShareResourceID(app, resourceID)
		if err != nil {
			return err
		}
		item, err := app.Client.FindShareLink(cmd.Context(), resourceType, id, args[0])
		if err != nil {
			return err
		}
		return app.OK(item, output.WithSummary(fmt.Sprintf("%s (%s)", item.Mode, item.Status)))
	}}
	addShareResourceFlags(cmd, &resourceType, &resourceID)
	return cmd
}
