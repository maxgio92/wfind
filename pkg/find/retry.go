package find

import (
	"time"

	"github.com/cenkalti/backoff"
)

var DefaultExponentialBackOffOptions = &ExponentialBackOffOptions{
	InitialInterval: 2 * time.Second,
	MaxInterval:     10 * time.Second,
	MaxElapsedTime:  5 * time.Minute,
}

type ExponentialBackOffOptions struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
	Clock           backoff.Clock
}

type ExponentialBackoffOption func(options *ExponentialBackOffOptions)

func WithInitialInterval(interval time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackOffOptions) {
		eb.InitialInterval = interval
	}
}
func WithMaxInterval(interval time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackOffOptions) {
		eb.MaxInterval = interval
	}
}

func WithMaxElapsedTime(elapsed time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackOffOptions) {
		eb.MaxElapsedTime = elapsed
	}
}

func WithClock(clock backoff.Clock) ExponentialBackoffOption {
	return func(eb *ExponentialBackOffOptions) {
		if clock != nil {
			eb.Clock = clock
		}
	}
}

func NewExponentialBackOffOptions(opts ...ExponentialBackoffOption) *ExponentialBackOffOptions {
	b := new(ExponentialBackOffOptions)

	for _, f := range opts {
		f(b)
	}

	return b
}
