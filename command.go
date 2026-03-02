package main

import (
	"context"
	"log/slog"
	"os"
	"syscall"

	"github.com/urfave/cli/v3"
)

func buildRootCommand() *cli.Command {
	return &cli.Command{
		Name:                  "core",
		Usage:                 "Get modules",
		EnableShellCompletion: true,
		Commands:              buildCommands(),
	}
}

type CommandGroup struct {
	BinaryFile
	CmdInfos []CmdInfo
}

func buildCommands() []*cli.Command {
	// get all binary names
	binaries := discoverBinaries()

	// fetch completion in JSON then group them with subcommand
	groups := make(map[string][]CommandGroup)
	for _, file := range binaries {
		cmdInfo, err := fetchCompletion(file.FullPath)
		if err != nil {
			slog.Error("get completion",
				slog.String("path", file.FullPath),
				slog.String("name", file.Name),
			)
			continue
		}
		cmdInfo.Name = file.Name

		// group with embed binary info (keep full binary path and command info list)
		for _, cmd := range cmdInfo.Commands {
			groups[cmd.Name] = append(groups[cmd.Name], CommandGroup{
				BinaryFile: file,
				CmdInfos:   cmd.Commands,
			})
		}
	}

	var cmds []*cli.Command
	for subCmd, group := range groups {
		cmd := mergeCommandGroup(subCmd, group)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func mergeCommandGroup(subCmd string, groups []CommandGroup) *cli.Command {
	var cmds []*cli.Command
	for _, cmd := range groups {
		for _, info := range cmd.CmdInfos {
			cmd := buildCommand(info, cmd.FullPath)
			cmds = append(cmds, cmd)
		}
	}

	return &cli.Command{
		Name:     subCmd,
		Commands: cmds,
	}
}

func buildCommand(info CmdInfo, binaryName string) *cli.Command {
	cmd := &cli.Command{
		Name:  info.Name,
		Usage: info.Usage,
		Flags: buildFlags(info.Flags),
	}

	// recursion create subcommand
	if len(info.Commands) > 0 {
		for _, subCmd := range info.Commands {
			cmd.Commands = append(cmd.Commands, buildCommand(subCmd, binaryName))
		}

		return cmd
	}

	// final node of subcommand
	cmd.Action = func(_ context.Context, c *cli.Command) error {
		// have to add argv[0] as binary name
		argv := append([]string{binaryName}, c.Root().Args().Slice()...)
		return syscall.Exec(binaryName, argv, os.Environ())
	}

	return cmd
}

func buildFlags(flags []string) []cli.Flag {
	// default from string first for POC only
	// production version must detect auto type
	// for boolean use case because boolean is TRUE means SET only
	// --enabled not --enabled=true
	cFlags := make([]cli.Flag, len(flags))
	for i, flag := range flags {
		cFlags[i] = &cli.StringFlag{Name: flag}
	}

	return cFlags
}
