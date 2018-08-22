package main

import (
	"context"
	"path/filepath"

	"github.com/docker/app/internal/content"
	"github.com/docker/cli/cli"
	"github.com/docker/docker/pkg/homedir"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	namespace string
	tag       string
}

func pushCmd() *cobra.Command {
	var opts pushOptions
	cmd := &cobra.Command{
		Use:   "push [<app-name>]",
		Short: "Push the application to a registry",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			/*
				app, err := packager.Extract(firstOrEmpty(args))
				if err != nil {
					return err
				}
				defer app.Cleanup()
			*/
			appStore, err := content.NewStore(filepath.Join(homedir.Get(), ".docker/app"))
			if err != nil {
				return err
			}
			ctx := context.Background()
			_, err = appStore.Push(ctx, args[0])
			return err
			// return packager.Push(app, opts.namespace, opts.tag)
		},
	}
	cmd.Flags().StringVar(&opts.namespace, "namespace", "", "namespace to use (default: namespace in metadata)")
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "tag to use (default: version in metadata)")
	return cmd
}
