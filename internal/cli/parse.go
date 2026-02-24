package cli

import (
	"errors"
	"flag"
	"strings"
)

type LaunchOptions struct {
	Tool      string
	Path      string
	ExtraArgs []string
}

func ParseLaunchArgs(args []string) (LaunchOptions, error) {
	mainArgs, extraArgs := splitPassthroughArgs(args)

	var opts LaunchOptions
	fs := flag.NewFlagSet("launch", flag.ContinueOnError)
	fs.StringVar(&opts.Tool, "tool", "", "tool to launch (codex|claude)")
	fs.StringVar(&opts.Path, "path", "", "file or directory path")

	if err := fs.Parse(mainArgs); err != nil {
		return LaunchOptions{}, err
	}

	opts.Tool = strings.TrimSpace(opts.Tool)
	opts.Path = strings.TrimSpace(opts.Path)
	if opts.Tool == "" {
		return LaunchOptions{}, errors.New("missing --tool")
	}
	if opts.Path == "" {
		return LaunchOptions{}, errors.New("missing --path")
	}

	opts.ExtraArgs = extraArgs
	return opts, nil
}

func splitPassthroughArgs(args []string) ([]string, []string) {
	for index, arg := range args {
		if arg == "--" {
			main := args[:index]
			extra := []string{}
			if index+1 < len(args) {
				extra = args[index+1:]
			}
			return main, extra
		}
	}
	return args, []string{}
}
