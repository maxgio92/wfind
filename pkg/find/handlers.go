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
	// Context timed out.
	case errors.Is(err, context.DeadlineExceeded):
		log.Println(err, "context deadline has exceeded")
	// Request has timed out.
	case os.IsTimeout(err):
		log.Println(err, "connection has timed out")
		if o.TimeoutRetryBackOff != nil {
			log.Println(err, "Will backoff...")
			retryWithExponentialBackoff(response.Request.Retry, o.TimeoutRetryBackOff)
		}
	// Connection has been reset (RST) by the peer.
	case errors.Is(err, unix.ECONNRESET):
		log.Println(err, "connection has been reset by peer")
		if o.ConnResetRetryBackOff != nil {
			log.Println(err, "Will backoff...")
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
			log.Println(err, "will retry...")
			continue
		}

		ticker.Stop()
		log.Println("retried with success")
		break
	}

	if err != nil {
		// Retry has failed.
		log.Println(err, "retry limit exhausted")
		return
	}

	// Retry is successful.
}
