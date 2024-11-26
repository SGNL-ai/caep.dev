package retry

import (
	"math"
	"math/rand/v2"
	"time"
)

// backoff implements exponential backoff with jitter
type backoff struct {
	current      time.Duration
	max          time.Duration
	multiplier   float64
	attemptCount int
}

func newBackoff(initial, max time.Duration, multiplier float64) *backoff {
	return &backoff{
		current:    initial,
		max:        max,
		multiplier: multiplier,
	}
}

// next returns the next backoff duration with jitter
func (b *backoff) next() time.Duration {
	defer func() {
		b.attemptCount++
	}()

	// Calculate the next backoff without jitter
	nextBackoff := float64(b.current) * math.Pow(b.multiplier, float64(b.attemptCount))

	if nextBackoff > float64(b.max) {
		nextBackoff = float64(b.max)
	}

	// Apply jitter (randomly between 80% and 100% of calculated duration)
	jitteredBackoff := time.Duration(nextBackoff * (0.8 + 0.2*rand.Float64()))

	return jitteredBackoff
}
