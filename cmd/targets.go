package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type targetConfig struct {
	proto   string
	address string
	display string
	label   string
}

type multiTargetFlag []string

func (m *multiTargetFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiTargetFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("target must not be empty")
	}
	*m = append(*m, value)
	return nil
}

func (m multiTargetFlag) Values() []string {
	return append([]string(nil), m...)
}

func readTargetsFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entries = append(entries, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func parseTargetSpec(spec, proto string) (targetConfig, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return targetConfig{}, fmt.Errorf("empty target specification")
	}

	var display string
	if eq := strings.Index(trimmed, "="); eq >= 0 {
		display = strings.TrimSpace(trimmed[:eq])
		trimmed = strings.TrimSpace(trimmed[eq+1:])
	}
	if trimmed == "" {
		return targetConfig{}, fmt.Errorf("missing host:port in target")
	}
	if !strings.Contains(trimmed, ":") {
		return targetConfig{}, fmt.Errorf("target must be in host:port format")
	}
	if display == "" {
		display = fmt.Sprintf("%s://%s", proto, trimmed)
	}
	return targetConfig{
		proto:   proto,
		address: trimmed,
		display: display,
		label:   padLabel(display, labelWidth),
	}, nil
}
