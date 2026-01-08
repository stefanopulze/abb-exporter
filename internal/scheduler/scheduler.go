package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		interval: interval,
		stopChan: make(chan struct{}),
		tickers:  make([]func(), 0),
	}
}

type Scheduler struct {
	ticker   *time.Ticker
	stopChan chan struct{}
	interval time.Duration
	tickers  []func()
}

func (s *Scheduler) Start() {
	slog.Info("Starting scheduler")
	s.ticker = time.NewTicker(s.interval)

	go func() {
		for {
			select {
			case <-s.stopChan:
				return

			case <-s.ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				s.process(ctx)
				cancel()
			}

		}
	}()
}

func (s *Scheduler) Stop() error {
	slog.Info("Stopping scheduler")
	s.ticker.Stop()
	s.stopChan <- struct{}{}
	return nil
}

func (s *Scheduler) process(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered panic scheduler: %v\n", r)
		}
	}()

	for _, t := range s.tickers {
		t()
	}
}

func (s *Scheduler) Tick(data func()) {
	s.tickers = append(s.tickers, data)
}
