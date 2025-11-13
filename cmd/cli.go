package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type appConfig struct {
	interval time.Duration
	timeout  time.Duration
	targets  []targetConfig
}

func parseCLI() appConfig {
	var targetFlags multiTargetFlag
	flag.Var(&targetFlags, "target", "host:port or name=host:port to probe (repeatable)")
	targetsFileFlag := flag.String("targets-file", "", "path to newline separated target definitions (same format as --target)")
	protoFlag := flag.String("proto", "tcp", "network protocol to use: tcp or udp")
	nameFlag := flag.String("name", "", "optional display name (applies when monitoring a single target)")
	intervalFlag := flag.Duration("interval", time.Second, "how often to measure latency")
	timeoutFlag := flag.Duration("timeout", 3*time.Second, "dial timeout per probe")
	flag.Parse()

	proto := strings.ToLower(strings.TrimSpace(*protoFlag))
	if proto != "tcp" && proto != "udp" {
		exitWithUsage("unsupported protocol: %s", proto)
	}

	timeout := *timeoutFlag
	if timeout <= 0 {
		exitWithUsage("timeout must be positive")
	}

	interval := *intervalFlag
	if interval <= 0 {
		exitWithUsage("interval must be positive")
	}

	targetEntries := targetFlags.Values()
	if file := strings.TrimSpace(*targetsFileFlag); file != "" {
		fileTargets, err := readTargetsFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read targets file: %v\n", err)
			os.Exit(1)
		}
		targetEntries = append(targetEntries, fileTargets...)
	}

	if len(targetEntries) == 0 {
		exitWithUsage("provide at least one --target or --targets-file entry")
	}

	targets := make([]targetConfig, 0, len(targetEntries))
	for _, spec := range targetEntries {
		cfg, err := parseTargetSpec(spec, proto)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid target %q: %v\n", spec, err)
			os.Exit(1)
		}
		targets = append(targets, cfg)
	}

	if len(targets) == 1 {
		if name := strings.TrimSpace(*nameFlag); name != "" {
			targets[0].display = name
			targets[0].label = padLabel(name, labelWidth)
		}
	} else if strings.TrimSpace(*nameFlag) != "" {
		fmt.Fprintln(os.Stderr, "--name is ignored when monitoring multiple targets")
	}

	return appConfig{
		interval: interval,
		timeout:  timeout,
		targets:  targets,
	}
}
