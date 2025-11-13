package main

import (
	"context"
	"fmt"
)

type renderer struct {
	blocks        [][]string
	lastLineCount int
}

func newRenderer(count int) *renderer {
	return &renderer{
		blocks: make([][]string, count),
	}
}

func (r *renderer) Init() {
	r.lastLineCount = r.printBlocks()
}

func (r *renderer) Run(ctx context.Context, updates <-chan statusUpdate) {
	if len(r.blocks) == 0 {
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
			if update.index >= 0 && update.index < len(r.blocks) {
				r.blocks[update.index] = update.lines
				r.renderAll()
			}
		}
	}
}

func (r *renderer) renderAll() {
	if r.lastLineCount > 0 {
		fmt.Printf("\033[%dA", r.lastLineCount)
	}
	r.lastLineCount = r.printBlocks()
}

func (r *renderer) printBlocks() int {
	total := 0
	for _, block := range r.blocks {
		if len(block) == 0 {
			fmt.Print("\r\n")
			total++
			continue
		}
		for _, line := range block {
			fmt.Printf("\r%s", line)
			fmt.Print("\033[K\n")
			total++
		}
	}
	return total
}
