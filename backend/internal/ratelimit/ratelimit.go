package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Limiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	rate    rate.Limit
	burst   int
	ttl     time.Duration
}

func New(r rate.Limit, burst int, ttl time.Duration) *Limiter {
	l := &Limiter{
		entries: map[string]*entry{},
		rate:    r,
		burst:   burst,
		ttl:     ttl,
	}

	go l.cleanupLoop()
	return l
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	e, exists := l.entries[key]
	if !exists {
		e = &entry{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.entries[key] = e
	}
	e.lastSeen = time.Now()
	lim := e.limiter
	l.mu.Unlock()

	return lim.Allow()
}

func (l *Limiter) cleanupLoop() {
	ticker := time.NewTicker(l.ttl)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		cutoff := time.Now().Add(-l.ttl)
		for key, e := range l.entries {
			if e.lastSeen.Before(cutoff) {
				delete(l.entries, key)
			}
		}
		l.mu.Unlock()
	}
}
