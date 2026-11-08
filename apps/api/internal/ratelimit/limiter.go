package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	rate   float64
	burst  float64
	mu     sync.Mutex
	now    func() time.Time
	bucket map[string]*bucket
}

type bucket struct {
	tokens float64
	last   time.Time
}

func New(perMinute int) *Limiter {
	limiter := &Limiter{
		now:    time.Now,
		bucket: make(map[string]*bucket),
	}
	if perMinute > 0 {
		limiter.rate = float64(perMinute) / 60
		limiter.burst = float64(perMinute)
	}
	return limiter
}

func (l *Limiter) Allow(key string) bool {
	if l.rate <= 0 || l.burst <= 0 {
		return true
	}
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.bucket[key]
	if !ok {
		if len(l.bucket) > 10000 {
			l.prune(now)
		}
		entry = &bucket{tokens: l.burst, last: now}
		l.bucket[key] = entry
	} else {
		elapsed := now.Sub(entry.last).Seconds()
		if elapsed > 0 {
			entry.tokens += elapsed * l.rate
			if entry.tokens > l.burst {
				entry.tokens = l.burst
			}
			entry.last = now
		}
	}
	if entry.tokens < 1 {
		return false
	}
	entry.tokens--
	return true
}

func (l *Limiter) prune(now time.Time) {
	for key, entry := range l.bucket {
		if now.Sub(entry.last) > 10*time.Minute {
			delete(l.bucket, key)
		}
	}
}
