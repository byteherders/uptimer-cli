package main

import (
	"context"
	"fmt"
)

type renderer struct {
	lines []string
}

func newRenderer(count int) *renderer {
	return &renderer{
		lines: make([]string, count),
	}
}

func (r *renderer) Init() {
	if len(r.lines) == 0 {
		return
	}
	for _, line := range r.lines {
		fmt.Println(line)
	}
}

func (r *renderer) Run(ctx context.Context, updates <-chan statusUpdate) {
	if len(r.lines) == 0 {
		<-ctx.Done()
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.index >= 0 && update.index < len(r.lines) {
				r.lines[update.index] = update.line
				r.renderAll()
			}
		}
	}
}

func (r *renderer) renderAll() {
	if len(r.lines) == 0 {
		return
	}
	fmt.Printf("\033[%dA", len(r.lines))
	for _, line := range r.lines {
		fmt.Printf("\r%s", line)
		fmt.Print("\033[K\n")
	}
}
