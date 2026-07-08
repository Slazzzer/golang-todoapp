package core_ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	max    int
	window time.Duration

	mu      sync.Mutex
	entries map[string]entry
}

type entry struct {
	count   int
	resetAt time.Time
}

func NewLimiter(max int, window time.Duration) *Limiter {
	return &Limiter{
		max:     max,
		window:  window,
		entries: make(map[string]entry),
	}
}

func (l *Limiter) Allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	current, ok := l.entries[key]
	if !ok || now.After(current.resetAt) {
		l.entries[key] = entry{
			count:   1,
			resetAt: now.Add(l.window),
		}

		return true
	}

	if current.count >= l.max {
		return false
	}

	current.count++
	l.entries[key] = current

	return true
}
