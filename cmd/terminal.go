package main

import (
	"os"
	"strconv"
)

func detectTerminalWidth() int {
	if w, err := terminalWidth(); err == nil && w > 0 {
		return w
	}
	if env := os.Getenv("COLUMNS"); env != "" {
		if n, err := strconv.Atoi(env); err == nil && n > 0 {
			return n
		}
	}
	return 0
}
