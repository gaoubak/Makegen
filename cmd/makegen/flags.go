package main

import "flag"

// ParseFlags parses command-line flags
func ParseFlags() {
	flag.BoolVar(verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(version_, "version", false, "Show version information")
	flag.BoolVar(help, "help", false, "Show help message")
	flag.Parse()
}
