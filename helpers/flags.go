package helpers

import (
	"flag"
	"os"

	"github.com/GDSC-Phenikaa/ctf-backend/globals"
	"github.com/fatih/color"
)

var (
	ShowHelp      = flag.Bool("help", false, "Show help message")
	ShowHelpShort = flag.Bool("h", false, "Show help message (shorthand)")
	Version       = flag.Bool("version", false, "Show version information")
	VersionShort  = flag.Bool("v", false, "Show version information")
)

func ParseFlags() {
	color.Set(color.FgMagenta, color.Bold)

	flag.Parse()

	if *ShowHelp || *ShowHelpShort {
		color.Set(color.FgMagenta, color.Bold)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *Version || *VersionShort {
		Help("GDSC CTF Backend API\nVersion: %s\nBuild Date: %s\nBuilt by: %s\n", globals.Version, globals.BuildDate, globals.BuildUser)
		os.Exit(0)
	}

	color.Unset()
}
