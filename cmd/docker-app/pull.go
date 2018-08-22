package main

import (
	"context"
	"path/filepath"

	"github.com/docker/app/internal/content"
	"github.com/docker/app/loader"
	"github.com/docker/cli/cli"
	"github.com/docker/docker/pkg/homedir"
	"github.com/spf13/cobra"
)

func pullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <repotag>",
		Short: "Pull an application from a registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// return packager.Pull(args[0])
			appStore, err := content.NewStore(filepath.Join(homedir.Get(), ".docker/app"))
			if err != nil {
				return err
			}
			ctx := context.Background()
			if err := appStore.Pull(ctx, args[0]); err != nil {
				return err
			}
			// If we want pull to also extract it!
			tar, err := appStore.Get(ctx, args[0])
			if err != nil {
				return err
			}
			app, err := loader.LoadFromTarReader(tar)
			if err != nil {
				return err
			}
			return app.Extract(".")
		},
	}
}
