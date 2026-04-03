package brickognize

import "time"

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
	if rps <= 0 {
		rps = 5
	}

	rl := &RateLimiter{
		tokens: make(chan struct{}, rps),
	}

	for i := 0; i < rps; i++ {
		rl.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rps))
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

func (r *RateLimiter) Wait() {
	<-r.tokens
}
