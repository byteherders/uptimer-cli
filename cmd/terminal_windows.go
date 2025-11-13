//go:build windows

package main

import "errors"

func terminalWidth() (int, error) {
	return 0, errors.New("terminal size detection not implemented on windows")
}
