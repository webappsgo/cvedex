package main

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestBuildDateFromEpoch(t *testing.T) {
	originalEpoch, originalDate := BuildEpoch, BuildDate
	t.Cleanup(func() {
		BuildEpoch, BuildDate = originalEpoch, originalDate
	})

	tests := []struct {
		name  string
		epoch string
		want  string
	}{
		{name: "valid epoch", epoch: "1700000000", want: "2023-11-14T22:13:20Z"},
		{name: "unset epoch", epoch: "0", want: "N/A"},
		{name: "negative epoch", epoch: "-1", want: "N/A"},
		{name: "invalid epoch", epoch: "invalid", want: "N/A"},
		{name: "empty epoch", epoch: "", want: "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildEpoch = tt.epoch
			BuildDate = "N/A"
			setBuildDateFromEpoch()
			if BuildDate != tt.want {
				t.Errorf("BuildDate = %q, want %q", BuildDate, tt.want)
			}
		})
	}
}

func TestMainStartsAndStopsOnSignal(t *testing.T) {
	t.Setenv("MODE", "development")
	t.Setenv("DEBUG", "false")

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	done := make(chan struct{})
	go func() {
		main()
		close(done)
	}()

	t.Cleanup(func() {
		os.Stdout = originalStdout
		_ = writer.Close()
		_ = reader.Close()
	})

	startup := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(reader).ReadString('\n')
		startup <- line
	}()

	select {
	case line := <-startup:
		if !strings.Contains(line, "starting in mode: development") {
			t.Fatalf("startup output %q does not report development mode", line)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("service did not report startup")
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("send SIGTERM: %v", err)
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("service did not stop after SIGTERM")
	}
}
