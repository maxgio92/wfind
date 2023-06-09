package find

import (
	"net/url"
	"regexp"
	"strings"

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

// NewFind returns a new Find object to find files over HTTP and HTTPS.
func NewFind(opts ...Option) *Options {
	o := &Options{}

	for _, f := range opts {
		f(o)
	}

	return o
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
