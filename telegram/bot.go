package telegram

import (
	"context"
	"sync"
	"time"

	"github.com/alireza0/s-ui/logger"
	"github.com/alireza0/s-ui/service"
)

type Bot struct {
	service *service.TelegramService
	mu      sync.Mutex
	cancel  context.CancelFunc
	running bool
	refresh chan struct{}
}

func NewBot(s *service.TelegramService) *Bot {
	if s == nil {
		s = service.SharedTelegramService()
	}
	bot := &Bot{
		service: s,
		refresh: make(chan struct{}, 1),
	}
	s.RegisterChangeListener(bot.TriggerRefresh)
	return bot
}

func (b *Bot) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.running {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	b.running = true
	go b.loop(ctx)
	logger.Info("telegram bot background worker started")
	return nil
}

func (b *Bot) loop(ctx context.Context) {
	defer func() {
		b.mu.Lock()
		b.running = false
		b.mu.Unlock()
	}()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	b.sync()
	for {
		select {
		case <-ctx.Done():
			logger.Info("telegram bot background worker stopped")
			return
		case <-ticker.C:
			b.sync()
		case <-b.refresh:
			b.sync()
		}
	}
}

func (b *Bot) sync() {
	state, err := b.service.GetAdminState()
	if err != nil {
		logger.Warning("telegram bot state refresh failed:", err)
		return
	}
	if state == nil || state.Config == nil {
		logger.Debug("telegram bot configuration is not yet initialized")
		return
	}
	if !state.Config.Enabled {
		logger.Debug("telegram bot is disabled; skipping synchronization")
		return
	}
	logger.Debugf("telegram bot synchronized with %d tariffs", len(state.Tariffs))
}

func (b *Bot) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cancel != nil {
		b.cancel()
		b.cancel = nil
	}
}

func (b *Bot) TriggerRefresh() {
	select {
	case b.refresh <- struct{}{}:
	default:
	}
}
