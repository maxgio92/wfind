package backoff

import (
	"time"

	"github.com/cenkalti/backoff"
)

type BackoffOption func(*backoff.ExponentialBackOff)

func WithInitialInterval(interval time.Duration) BackoffOption {
	return func(eb *backoff.ExponentialBackOff) {
		eb.InitialInterval = interval
	}
}
func WithMaxInterval(interval time.Duration) BackoffOption {
	return func(eb *backoff.ExponentialBackOff) {
		eb.MaxInterval = interval
	}
}

func WithMaxElapsedTime(elapsed time.Duration) BackoffOption {
	return func(eb *backoff.ExponentialBackOff) {
		eb.MaxElapsedTime = elapsed
	}
}

func WithClock(clock backoff.Clock) BackoffOption {
	return func(eb *backoff.ExponentialBackOff) {
		if clock != nil {
			eb.Clock = clock
		}
	}
}

func NewExponentialBackOff(opts ...BackoffOption) *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
		Clock:               backoff.SystemClock,
	}

	for _, f := range opts {
		f(b)
	}
	b.Reset()

	return b
}

// Clock is an interface that returns current time for BackOff.
type Clock interface {
	Now() time.Time
}
