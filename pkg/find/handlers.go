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
	"context"
	"log"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	utils "github.com/maxgio92/wfind/internal/backoff"
)

// handleError handles an error received making a colly.Request.
// It accepts a colly.Response and the error.
func (o *Options) handleError(response *colly.Response, err error) {
	switch {
	// Context deadline is passed.
	case errors.Is(err, context.DeadlineExceeded):
		if o.ContextDeadlineRetryBackOff != nil {
			retryWithExponentialBackoff(response.Request.Retry, o.ContextDeadlineRetryBackOff)
		}
	// Request has timed out.
	case os.IsTimeout(err):
		if o.TimeoutRetryBackOff != nil {
			retryWithExponentialBackoff(response.Request.Retry, o.TimeoutRetryBackOff)
		}
	// Connection has been reset (RST) by the peer.
	case errors.Is(err, unix.ECONNRESET):
		if o.ConnResetRetryBackOff != nil {
			retryWithExponentialBackoff(response.Request.Retry, o.ConnResetRetryBackOff)
		}
	// Other failures.
	default:
		log.Printf("error: %v\n", err)
	}
}

// retryWithExtponentialBackoff retries with an exponential backoff a function.
// Exponential backoff can be tuned with options accepted as arguments to the function.
func retryWithExponentialBackoff(retryF func() error, opts *ExponentialBackOffOptions) {
	ticker := backoff.NewTicker(
		utils.NewExponentialBackOff(
			utils.WithClock(opts.Clock),
			utils.WithInitialInterval(opts.InitialInterval),
			utils.WithMaxInterval(opts.MaxInterval),
			utils.WithMaxElapsedTime(opts.MaxElapsedTime),
		),
	)

	var err error

	// Ticks will continue to arrive when the previous retryF is still running,
	// so operations that take a while to fail could run in quick succession.
	for range ticker.C {
		if err = retryF(); err != nil {
			// Retry.
			continue
		}

		ticker.Stop()
		break
	}

	if err != nil {
		// Retry has failed.
		return
	}
}
