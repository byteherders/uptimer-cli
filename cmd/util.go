package main

import (
	"flag"
	"fmt"
	"os"
)

func exitWithUsage(format string, args ...interface{}) {
	if len(format) > 0 {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
	flag.Usage()
	os.Exit(2)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
