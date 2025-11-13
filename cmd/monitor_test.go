package main

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestVerifyConnectivitySuccess(t *testing.T) {
	restore := stubDial(func(network, address string, timeout time.Duration) (net.Conn, error) {
		c1, c2 := net.Pipe()
		go c2.Close()
		return c1, nil
	})
	defer restore()

	if err := verifyConnectivity("tcp", "example.com:80", time.Second); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestVerifyConnectivityFailure(t *testing.T) {
	restore := stubDial(func(network, address string, timeout time.Duration) (net.Conn, error) {
		return nil, errors.New("boom")
	})
	defer restore()

	if err := verifyConnectivity("tcp", "example.com:80", time.Second); err == nil {
		t.Fatalf("expected failure")
	}
}

func TestMeasureLatency(t *testing.T) {
	restore := stubDial(func(network, address string, timeout time.Duration) (net.Conn, error) {
		c1, c2 := net.Pipe()
		go c2.Close()
		time.Sleep(5 * time.Millisecond)
		return c1, nil
	})
	defer restore()

	latency, err := measureLatency("tcp", "example.com:80", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latency <= 0 {
		t.Fatalf("expected positive latency, got %v", latency)
	}
}

func stubDial(fn func(string, string, time.Duration) (net.Conn, error)) func() {
	original := dialTimeout
	dialTimeout = fn
	return func() {
		dialTimeout = original
	}
}
