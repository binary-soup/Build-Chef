package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/binary-soup/bchef/cmd"
	"github.com/binary-soup/bchef/common"
	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/style"
)

const (
	VERSION = "v0.1.0-alpha"
)

var cmds = []cmd.Command{
	cmd.NewCleanCommand(),
	cmd.NewCookCommand(),
	cmd.NewMakeCommand(),
	cmd.NewReviewCommand(),
}

func main() {
	if handleFlags() {
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("no command given")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		style.Info.Println("[USING DEFAULT CONFIG]")
	} else {
		style.Info.Println("[CONFIG LOADED]")
	}

	if err = runCommand(os.Args[1], cfg, os.Args[2:]); err != nil {
		common.PrintError(err)
	}
}

func handleFlags() bool {
	version := flag.Bool("v", false, "version info")
	list := flag.Bool("ls", false, "list commands")

	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Println(styledVersion())
		return true
	}

	if *list {
		printCommands()
		return true
	}

	return false
}

func usage() {
	fmt.Printf("%s (%s) ~ Build a c++ project using recipe files\n%s\n",
		style.BoldInfoV2.String("Build Chef"), styledVersion(), style.BoldFileV2.String("Options:"))

	flag.PrintDefaults()
}

func styledVersion() string {
	return style.File.String(VERSION)
}

func printCommands() {
	for _, cmd := range cmds {
		cmd.PrintUsage()
	}
}

func runCommand(name string, cfg config.Config, args []string) error {
	for _, cmd := range cmds {
		if cmd.GetName() == name {
			return cmd.Run(cfg, args)
		}
	}
	return fmt.Errorf("unknown command \"%s\"", name)
}
