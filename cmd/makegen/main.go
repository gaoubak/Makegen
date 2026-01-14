package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gaoubak/Makegen/internal/app"
	"github.com/gaoubak/Makegen/internal/utils"
)

var (
	version  = "1.0.0"
	verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	version_ = flag.Bool("version", false, "Show version")
	help     = flag.Bool("help", false, "Show help")
)

func main() {
	flag.Parse()

	if *version_ {
		fmt.Printf("makegen version %s\n", version)
		os.Exit(0)
	}

	if *help {
		showHelp()
		os.Exit(0)
	}

	// Initialize logger
	logger := utils.NewLogger(*verbose)

	// Get working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	// Create and run application
	application := app.NewApp(logger, workDir)
	if err := application.Run(); err != nil {
		logger.Error("Application error: %v", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`🔨 Makefile Generator - Interactive Makefile Creation

Usage:
  makegen [flags]

Flags:
  -verbose    Enable verbose output
  -version    Show version
  -help       Show this help message

Examples:
  makegen                  Run interactive generator
  makegen -verbose         Run with debug output
  makegen -version         Show version

For more information, visit: https://github.com/yourusername/makegen
`)
}
