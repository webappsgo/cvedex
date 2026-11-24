package service

import (
	"bufio"
	"io"
	"os"
	"testing"
	"time"

	"github.com/casapps/cvedex/src/config"
)

func TestNewAndRunUntilShutdown(t *testing.T) {
	cfg := &config.Config{
		Version:  "test-version",
		CommitID: "test-commit",
		Mode:     "test-mode",
	}
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if svc.cfg != cfg {
		t.Fatal("New() did not retain the supplied configuration")
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = oldStdout
		_ = writer.Close()
		_ = reader.Close()
	})

	lineCh := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(reader).ReadString('\n')
		lineCh <- line
	}()
	runCh := make(chan error, 1)
	go func() {
		runCh <- svc.Run()
	}()

	select {
	case got := <-lineCh:
		if want := "cvedex test-version (test-commit) starting in mode: test-mode\n"; got != want {
			t.Errorf("Run() startup output = %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		svc.Shutdown()
		t.Fatal("Run() did not write its startup message")
	}

	select {
	case err := <-runCh:
		t.Fatalf("Run() returned before shutdown: %v", err)
	default:
	}

	svc.Shutdown()
	svc.Shutdown()
	select {
	case err := <-runCh:
		if err != nil {
			t.Errorf("Run() error after shutdown = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run() did not return after shutdown")
	}

	if err := writer.Close(); err != nil {
		t.Errorf("writer.Close() error = %v", err)
	}
	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Errorf("draining startup output: %v", err)
	}
}
