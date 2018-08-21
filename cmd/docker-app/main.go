package main

import (
	"context"
	"os"

	"github.com/docker/app/internal/com"

	"github.com/docker/cli/cli/command"
	"github.com/sirupsen/logrus"
)

func main() {
	_, streams, session, err := com.ConnectToFront(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
	// Set terminal emulation based on platform as required.

	logrus.SetOutput(streams.Err)

	dockerCli := command.NewDockerCli(streams.In, streams.Out, streams.Err, false)
	cmd := newRootCmd(dockerCli)
	cmd.SetOutput(streams.Err)
	err = cmd.Execute()
	com.Shutdown(context.Background(), session, streams)
	if err != nil {
		os.Exit(1)
	}
}
