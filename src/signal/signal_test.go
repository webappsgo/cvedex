package signal

import (
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

type recordingShutdowner struct {
	calls atomic.Int32
	done  chan struct{}
}

func (s *recordingShutdowner) Shutdown() {
	s.calls.Add(1)
	s.done <- struct{}{}
}

func TestHandleShutsDownOnSupportedSignals(t *testing.T) {
	tests := []struct {
		name string
		sig  os.Signal
	}{
		{name: "terminate", sig: syscall.SIGTERM},
		{name: "interrupt", sig: syscall.SIGINT},
		{name: "hangup", sig: syscall.SIGHUP},
	}
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("os.FindProcess() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdowner := &recordingShutdowner{done: make(chan struct{}, 1)}
			Handle(shutdowner)

			if err := process.Signal(tt.sig); err != nil {
				t.Fatalf("sending %v: %v", tt.sig, err)
			}
			select {
			case <-shutdowner.done:
				if got := shutdowner.calls.Load(); got != 1 {
					t.Errorf("Shutdown() calls = %d, want 1", got)
				}
			case <-time.After(3 * time.Second):
				t.Fatalf("Shutdown() was not called after %v", tt.sig)
			}
		})
	}
}
