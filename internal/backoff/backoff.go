/*
Copyright © 2023 maxgio92 me@maxgio.me

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
