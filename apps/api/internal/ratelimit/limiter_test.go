package ratelimit

import "testing"

func TestBurstThenDeny(t *testing.T) {
	limiter := New(2)
	if !limiter.Allow("a") || !limiter.Allow("a") {
		t.Fatal("first two requests should pass")
	}
	if limiter.Allow("a") {
		t.Fatal("third request should be denied")
	}
	if !limiter.Allow("b") {
		t.Fatal("different client should have its own bucket")
	}
}

func TestDisabledLimiter(t *testing.T) {
	limiter := New(0)
	for range 10 {
		if !limiter.Allow("a") {
			t.Fatal("disabled limiter must allow all requests")
		}
	}
}
