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
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/maxgio92/wfind/internal/network"
	"github.com/pkg/errors"
)

// Result represents the output of the Find job.
type Result struct {
	// BaseNames are the path base of the files found.
	BaseNames []string

	// URLs are the universal resource location of the files found.
	URLs []string
}

// Options represents the options for the Find job.
type Options struct {
	// SeedURLs are the URLs used as root URLs from which for the find's web scraping.
	SeedURLs []string

	// FilenameRegexp is a regular expression for which a pattern should match the file names in the Result.
	FilenameRegexp string

	// FileType is the file type for which the Find job examines the web hierarchy.
	FileType string

	// Recursive enables the Find job to examine files referenced to by the seeds files recursively.
	Recursive bool

	// Verbose enables the Find job verbosity printing every visited URL.
	Verbose bool

	// Async represetns the option to scrape with multiple asynchronous coroutines.
	Async bool

	// ClientTransport represents the Transport used for the HTTP client.
	ClientTransport http.RoundTripper

	// MaxBodySize is the limit in bytes of each of the retrieved response body.
	MaxBodySize int

	// ContextDeadlineRetryBackOff controls the error handling on responses.
	// If not nil, when the request context deadline exceeds, the request
	// is retried with an exponential backoff interval.
	ContextDeadlineRetryBackOff *ExponentialBackOffOptions

	// ConnResetRetryBackOff controls the error handling on responses.
	// If not nil, when the connection is reset by the peer (TCP RST), the request
	// is retried with an exponential backoff interval.
	ConnResetRetryBackOff *ExponentialBackOffOptions

	// TimeoutRetryBackOff controls the error handling on responses.
	// If not nil, when the connection times out (based on client timeout), the request
	// is retried with an exponential backoff interval.
	TimeoutRetryBackOff *ExponentialBackOffOptions
}

type Option func(opts *Options)

func WithSeedURLs(seedURLs []string) Option {
	return func(opts *Options) {
		opts.SeedURLs = seedURLs
	}
}

func WithFilenameRegexp(filenameRegexp string) Option {
	return func(opts *Options) {
		opts.FilenameRegexp = filenameRegexp
	}
}

func WithFileType(fileType string) Option {
	return func(opts *Options) {
		opts.FileType = fileType
	}
}

func WithRecursive(recursive bool) Option {
	return func(opts *Options) {
		opts.Recursive = recursive
	}
}

func WithVerbosity(verbosity bool) Option {
	return func(opts *Options) {
		opts.Verbose = verbosity
	}
}

func WithAsync(async bool) Option {
	return func(opts *Options) {
		opts.Async = async
	}
}

func WithClientTransport(transport http.RoundTripper) Option {
	return func(opts *Options) {
		opts.ClientTransport = transport
	}
}

func WithMaxBodySize(maxBodySize int) Option {
	return func(opts *Options) {
		opts.MaxBodySize = maxBodySize
	}
}

func WithContextDeadlineRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.ContextDeadlineRetryBackOff = backoff
	}
}

func WithConnResetRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.ConnResetRetryBackOff = backoff
	}
}

func WithConnTimeoutRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.TimeoutRetryBackOff = backoff
	}
}

// NewFind returns a new Find object to find files over HTTP and HTTPS.
func NewFind(opts ...Option) *Options {
	o := &Options{}

	for _, f := range opts {
		f(o)
	}

	o.init()

	return o
}

// Validate validates the Find job options and returns an error.
func (o *Options) init() {
	if o.ClientTransport == nil {
		o.ClientTransport = network.DefaultClientTransport
	}
	if o.MaxBodySize == 0 {
		// Set max body size to 100 KB.
		o.MaxBodySize = 100 * 1024
	}
}

// Validate validates the Find job options and returns an error.
func (o *Options) Validate() error {
	// Validate seed URLs.
	if len(o.SeedURLs) == 0 {
		return errors.New("no seed URLs specified")
	}

	for k, v := range o.SeedURLs {
		_, err := url.Parse(v)
		if err != nil {
			return errors.New("a seed URL is not a valid URL")
		}

		if !strings.HasSuffix(v, "/") {
			o.SeedURLs[k] = v + "/"
		}
	}

	// Validate filename regular expression.
	if o.FilenameRegexp == "" {
		return errors.New("no filename regular expression specified")
	}

	if _, err := regexp.Compile(o.FilenameRegexp); err != nil {
		return errors.Wrap(err, "error validating the file name expression")
	}

	// Validate file type.
	if o.FileType == "" {
		o.FileType = FileTypeReg
	} else if o.FileType != FileTypeReg && o.FileType != FileTypeDir {
		return errors.New("file type not supported")
	}

	o.sanitize()

	return nil
}

func (o *Options) sanitize() {
	if strings.HasPrefix(o.FilenameRegexp, "^") && !strings.HasPrefix(o.FilenameRegexp, "^./") && !strings.HasPrefix(o.FilenameRegexp, `^(\./)?`) {
		o.FilenameRegexp = strings.Replace(o.FilenameRegexp, "^", `^(\./)?`, 1)
	}

	if o.FileType == FileTypeDir {
		if strings.HasSuffix(o.FilenameRegexp, "$") && !strings.HasSuffix(o.FilenameRegexp, "/$") {
			last := strings.LastIndex(o.FilenameRegexp, "$")
			o.FilenameRegexp = o.FilenameRegexp[:last] + strings.Replace(o.FilenameRegexp[last:], "$", "/?$", 1)
		}
	}
}

func (o *Options) Find() (*Result, error) {
	if err := o.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating find options")
	}

	switch o.FileType {
	case FileTypeReg:
		return o.crawlFiles()
	case FileTypeDir:
		return o.crawlFolders()
	default:
		return o.crawlFiles()
	}
}
