package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"dir-agent/internal/cli"
	"dir-agent/internal/config"
	"dir-agent/internal/launcher"
	"dir-agent/internal/resources"
	"dir-agent/internal/terminal"
)

var version = "dev"

func main() {
	exitCode := run(os.Args[1:])
	os.Exit(exitCode)
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 2
	}

	switch args[0] {
	case "launch":
		return runLaunch(args[1:])
	case "doctor":
		return runDoctor()
	case "version":
		fmt.Println(version)
		return 0
	case "install-assets":
		return runInstallAssets()
	case "uninstall-assets":
		return runUninstallAssets(args[1:])
	case "path":
		return runPath(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] unknown command: %s\n", args[0])
		printUsage()
		return 2
	}
}

func runLaunch(args []string) int {
	opts, err := cli.ParseLaunchArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] %v\n", err)
		return 2
	}

	cfgPath, err := config.EnsureConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] resolve config path: %v\n", err)
		return 2
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] load config: %v\n", err)
		return 2
	}

	req := launcher.LaunchRequest{
		ToolName:  opts.Tool,
		InputPath: opts.Path,
		ExtraArgs: opts.ExtraArgs,
		Config:    cfg,
	}
	if err := launcher.Launch(req); err != nil {
		switch {
		case errors.Is(err, launcher.ErrToolNotFound):
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] tool missing in PATH: %v\n", err)
			return 3
		case errors.Is(err, launcher.ErrPathNotAccessible):
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] path not accessible: %v\n", err)
			return 5
		case errors.Is(err, terminal.ErrNoTerminalFound):
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] no supported terminal found. set [terminals].preferred in config\n")
			return 4
		default:
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] launch failed: %v\n", err)
			return 2
		}
	}

	fmt.Printf("[diragent] Launching %s in %s\n", opts.Tool, opts.Path)
	return 0
}

func runDoctor() int {
	cfgPath, err := config.EnsureConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] config path: %v\n", err)
		return 2
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] load config: %v\n", err)
		return 2
	}

	fmt.Printf("Config: %s\n", cfgPath)
	fmt.Printf("Preferred terminal: %s\n", cfg.Terminals.Preferred)

	codexState := "missing"
	if _, err := terminal.FindExecutable(cfg.Tools.Codex.Command); err == nil {
		codexState = "ok"
	}
	claudeState := "missing"
	if _, err := terminal.FindExecutable(cfg.Tools.Claude.Command); err == nil {
		claudeState = "ok"
	}

	fmt.Printf("codex: %s\n", codexState)
	fmt.Printf("claude: %s\n", claudeState)
	return 0
}

func runInstallAssets() int {
	result, err := resources.Install()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] install assets: %v\n", err)
		return 2
	}

	output, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] marshal result: %v\n", err)
		return 2
	}
	fmt.Println(string(output))
	return 0
}

func runUninstallAssets(args []string) int {
	fs := flag.NewFlagSet("uninstall-assets", flag.ContinueOnError)
	removeConfig := fs.Bool("remove-config", false, "remove config file")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if err := resources.Uninstall(*removeConfig); err != nil {
		fmt.Fprintf(os.Stderr, "[diragent][ERROR] uninstall assets: %v\n", err)
		return 2
	}
	return 0
}

func runPath(args []string) int {
	fs := flag.NewFlagSet("path", flag.ContinueOnError)
	kind := fs.String("kind", "", "one of: data, config")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	switch *kind {
	case "data":
		value, err := config.DataPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] data path: %v\n", err)
			return 2
		}
		fmt.Println(value)
		return 0
	case "config":
		value, err := config.ActiveConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[diragent][ERROR] config path: %v\n", err)
			return 2
		}
		fmt.Println(value)
		return 0
	default:
		fmt.Fprintln(os.Stderr, "[diragent][ERROR] path --kind must be data or config")
		return 2
	}
}

func printUsage() {
	fmt.Println(`diragent commands:
  launch --tool codex|claude --path <file-or-dir> [-- <extra args...>]
  doctor
  version
  path --kind data|config
  install-assets
  uninstall-assets [--remove-config]`)
}
