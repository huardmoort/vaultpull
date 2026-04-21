package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/runner"
)

func main() {
	cfgPath := flag.String("config", "vaultpull.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	results := runner.Run(cfg)

	os.Exit(printResults(results))
}

// printResults writes a summary of each run result to stdout or stderr and
// returns a non-zero exit code if any result contained an error.
func printResults(results []runner.Result) int {
	exitCode := 0
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "[FAIL] %s: %v\n", r.EnvFile, r.Err)
			exitCode = 1
		} else {
			fmt.Printf("[OK]   %s (%d secret(s) written)\n", r.EnvFile, r.Written)
		}
	}
	return exitCode
}
