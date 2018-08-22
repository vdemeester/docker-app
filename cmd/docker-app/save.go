package main

import (
	"context"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/docker/app/internal"
	"github.com/docker/app/internal/content"
	"github.com/docker/app/internal/packager"
	"github.com/docker/app/types/metadata"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/pkg/homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type saveOptions struct {
	namespace string
	tag       string
}

func saveCmd(dockerCli command.Cli) *cobra.Command {
	var opts saveOptions
	cmd := &cobra.Command{
		Use:   "save [<app-name>]",
		Short: "Save the application as an image to the docker daemon(in preparation for push)",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := packager.Extract(firstOrEmpty(args))
			if err != nil {
				return err
			}
			defer app.Cleanup()
			appStore, err := content.NewStore(filepath.Join(homedir.Get(), ".docker/app"))
			if err != nil {
				return err
			}
			ctx := context.Background()
			namespace := opts.namespace
			tag := opts.tag
			var meta metadata.AppMetadata
			err = yaml.Unmarshal(app.Metadata(), &meta)
			if err != nil {
				return errors.Wrap(err, "failed to parse application metadata")
			}
			if tag == "" {
				tag = meta.Version
			}
			if namespace == "" {
				namespace = meta.Namespace
			}
			if namespace != "" && !strings.HasSuffix(namespace, "/") {
				namespace += "/"
			}
			imageName := namespace + internal.AppNameFromDir(app.Name) + internal.AppExtension + ":" + tag
			r, err := app.TarReader()
			if err != nil {
				return err
			}
			defer r.Close()
			return appStore.Load(ctx, imageName, r)
			/*
				imageName, err := packager.Save(app, opts.namespace, opts.tag)
				if imageName != "" && err == nil {
					fmt.Fprintf(dockerCli.Out(), "Saved application as image: %s\n", imageName)
				}
				return err
			*/
		},
	}
	cmd.Flags().StringVar(&opts.namespace, "namespace", "", "namespace to use (default: namespace in metadata)")
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "tag to use (default: version in metadata)")
	return cmd
}
