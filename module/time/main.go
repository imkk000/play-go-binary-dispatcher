package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

var rootCmd = &cli.Command{
	Name:  "root",
	Usage: "Go Time POC",
	Commands: []*cli.Command{
		subGetCmd,
		subCompCmd,
	},
	CommandNotFound: func(_ context.Context, c *cli.Command, s string) {
		fmt.Println(c.Args(), s)
	},
}

var subGetCmd = &cli.Command{
	Name:  "get",
	Usage: "Get value",
	Commands: []*cli.Command{
		{
			Name:  "time",
			Usage: "Get local time",
			Commands: []*cli.Command{
				{
					Name:  "now",
					Usage: "Get current time now",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "TZ",
						},
					},
					Action: execAction,
				},
			},
		},
	},
}

var subCompCmd = &cli.Command{
	Name:  "completion",
	Usage: "Get shell completion",
	Commands: []*cli.Command{
		{
			Name:  "json",
			Usage: "Get completion in json format",
			Action: func(_ context.Context, c *cli.Command) error {
				info := walkCmd(c.Root())
				data, err := json.Marshal(info)
				if err != nil {
					return err
				}
				fmt.Println(string(data))

				return nil
			},
		},
	},
}

type CmdInfo struct {
	Name     string    `json:"name"`
	Usage    string    `json:"usage"`
	Commands []CmdInfo `json:"commands,omitempty"`
	Flags    []string  `json:"flags,omitempty"`
}

type SkipCommand map[string]bool

var skipCommand = SkipCommand{
	"help":       true,
	"h":          true,
	"version":    true,
	"completion": true,
}

func (c SkipCommand) Is(cmd string) bool {
	return c[cmd]
}

func walkCmd(c *cli.Command) CmdInfo {
	info := CmdInfo{
		Name:  c.Name,
		Usage: c.Usage,
	}

	for _, cmd := range c.Commands {
		if skipCommand.Is(cmd.Name) {
			continue
		}

		info.Commands = append(info.Commands, walkCmd(cmd))
	}
	for _, flag := range c.Flags {
		var names []string
		for _, name := range flag.Names() {
			if skipCommand.Is(name) {
				continue
			}
			names = append(names, name)
		}
		info.Flags = append(info.Flags, names...)
	}

	return info
}

func execAction(_ context.Context, c *cli.Command) error {
	slog.Info("exec "+c.Name,
		slog.Any("args", c.Args()),
		slog.Any("flags", c.FlagNames()),
	)

	return nil
}

func main() {
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("run command", slog.Any("err", err))
		os.Exit(1)
	}
}
