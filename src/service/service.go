// Package service implements the top-level service lifecycle for cvedex.
package service

import (
	"context"
	"fmt"

	"github.com/casapps/cvedex/src/config"
)

// Service is the top-level cvedex service — owns the DNS server, HTTP gateway,
// scheduler, and graceful shutdown coordination.
type Service struct {
	cfg *config.Config
	ctx context.Context
	cancel context.CancelFunc
}

// New constructs a Service with the provided configuration.
func New(cfg *config.Config) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Run starts cvedex and blocks until shutdown is requested.
func (s *Service) Run() error {
	fmt.Printf("cvedex %s (%s) starting in mode: %s\n",
		s.cfg.Version, s.cfg.CommitID, s.cfg.Mode)

	<-s.ctx.Done()
	return nil
}

// Shutdown requests a graceful stop of all service components.
func (s *Service) Shutdown() {
	s.cancel()
}
