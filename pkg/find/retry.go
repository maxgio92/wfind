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
