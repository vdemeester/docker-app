package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/app/internal/content"
	"github.com/docker/cli/cli"
	"github.com/docker/docker/pkg/homedir"
	"github.com/spf13/cobra"
)

type listOptions struct {
	quiet bool
}

func lsCmd() *cobra.Command {
	var opts listOptions
	cmd := &cobra.Command{
		Use:   "ls [<app-name>:[<tag>]]",
		Short: "List applications.",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//			return image.List(firstOrEmpty(args), opts.quiet)

			appStore, err := content.NewStore(filepath.Join(homedir.Get(), ".docker/app"))
			if err != nil {
				return err
			}
			ctx := context.Background()
			apps, err := appStore.List(ctx)
			if err != nil {
				return err
			}
			for _, app := range apps {
				fmt.Println(app.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}
