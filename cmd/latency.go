package main

import "time"

type latencyBuffer struct {
	values []time.Duration
	sum    time.Duration
	next   int
	full   bool
}

type latencyStats struct {
	avg   time.Duration
	min   time.Duration
	max   time.Duration
	count int
}

func newLatencyBuffer(size int) *latencyBuffer {
	if size <= 0 {
		size = 1
	}
	return &latencyBuffer{
		values: make([]time.Duration, size),
	}
}

func (b *latencyBuffer) Add(v time.Duration) {
	if len(b.values) == 0 {
		return
	}
	if b.full {
		b.sum -= b.values[b.next]
	}
	b.values[b.next] = v
	b.sum += v
	b.next++
	if b.next >= len(b.values) {
		b.next = 0
		b.full = true
	}
}

func (b *latencyBuffer) Stats() latencyStats {
	count := b.count()
	if count == 0 {
		return latencyStats{}
	}
	avg := time.Duration(int64(b.sum) / int64(count))
	min := time.Duration(0)
	max := time.Duration(0)
	first := true

	iterateValues(b, func(v time.Duration) {
		if first {
			min = v
			max = v
			first = false
			return
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	})

	return latencyStats{
		avg:   avg,
		min:   min,
		max:   max,
		count: count,
	}
}

func (b *latencyBuffer) Values() []time.Duration {
	count := b.count()
	result := make([]time.Duration, 0, count)
	if count == 0 {
		return result
	}
	if b.full {
		result = append(result, b.values[b.next:]...)
		result = append(result, b.values[:b.next]...)
	} else {
		result = append(result, b.values[:b.next]...)
	}
	return result
}

func (b *latencyBuffer) count() int {
	if b.full {
		return len(b.values)
	}
	return b.next
}

func iterateValues(b *latencyBuffer, fn func(time.Duration)) {
	if b == nil || fn == nil {
		return
	}
	if b.full {
		for i := 0; i < len(b.values); i++ {
			idx := (b.next + i) % len(b.values)
			fn(b.values[idx])
		}
	} else {
		for i := 0; i < b.next; i++ {
			fn(b.values[i])
		}
	}
}
