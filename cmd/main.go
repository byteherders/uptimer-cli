package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

const (
	graphWidth     = 48
	labelWidth     = 28
	historySamples = 256
)

func main() {
	cfg := parseCLI()

	for _, target := range cfg.targets {
		if err := verifyConnectivity(target.proto, target.address, cfg.timeout); err != nil {
			fmt.Fprintf(os.Stderr, "failed to reach %s (%s): %v\n", target.display, target.address, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Monitoring %d target(s) every %s (Ctrl+C to stop)\n", len(cfg.targets), cfg.interval)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	runMonitors(ctx, cfg)

	fmt.Println("\nExiting...")
}
