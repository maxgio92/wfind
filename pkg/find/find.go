package find

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
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
	if strings.HasPrefix(o.FilenameRegexp, "^") && !strings.HasPrefix(o.FilenameRegexp, "^./") {
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

// crawlFiles returns a list of file names found from the seed URL, filtered by file name regex.
//
//nolint:funlen,cyclop
func (o *Options) crawlFiles() (*Result, error) {
	seeds := []*url.URL{}

	err := o.Validate()
	if err != nil {
		return nil, err
	}

	for _, v := range o.SeedURLs {
		u, _ := url.Parse(v)

		seeds = append(seeds, u)
	}

	var files, urls []string

	folderPattern := regexp.MustCompile(folderRegex)

	exactFilePattern := regexp.MustCompile(o.FilenameRegexp)

	fileRegex := strings.TrimPrefix(o.FilenameRegexp, "^")
	filePattern := regexp.MustCompile(fileRegex)

	allowedDomains := getHostnamesFromURLs(seeds)

	// Create the collector settings
	coOptions := []func(*colly.Collector){
		colly.AllowedDomains(allowedDomains...),
		colly.Async(false),
	}

	if o.Verbose {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)

	// Add the callback to Visit the linked resource, for each HTML element found
	co.OnHTML(HTMLTagLink, func(e *colly.HTMLElement) {
		link := e.Attr(HTMLAttrRef)

		// Do not traverse the hierarchy in reverse order.
		if o.Recursive && !(strings.Contains(link, UpDir)) && link != RootDir {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Add the analysis callback to find file URLs, for each Visit call
	co.OnRequest(func(r *colly.Request) {
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// If the URL is not of a folder.
		if len(folderMatch) == 0 {
			fileMatch := filePattern.FindStringSubmatch(r.URL.String())

			// If the URL is of a file.
			if len(fileMatch) > 0 {
				fileName := path.Base(r.URL.String())
				fileNameMatch := exactFilePattern.FindStringSubmatch(fileName)

				// If the URL matches the file filter regex.
				if len(fileNameMatch) > 0 {
					files = append(files, fileName)
					urls = append(urls, r.URL.String())
				}
			}
			// Otherwise abort the request.
			r.Abort()
		}
	})

	// Visit each root folder.
	for _, seedURL := range seeds {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error scraping file with URL %seedURLs", seedURL.String()))
		}
	}

	return &Result{BaseNames: files, URLs: urls}, nil
}

// crawlFolders returns a list of folder names found from each seed URL, filtered by folder name regex.
//
//nolint:funlen,cyclop
func (o *Options) crawlFolders() (*Result, error) {
	seeds := []*url.URL{}

	err := o.Validate()
	if err != nil {
		return nil, err
	}

	for _, v := range o.SeedURLs {
		u, _ := url.Parse(v)

		seeds = append(seeds, u)
	}

	var folders, urls []string

	folderPattern := regexp.MustCompile(folderRegex)

	exactFolderPattern := regexp.MustCompile(o.FilenameRegexp)

	allowedDomains := getHostnamesFromURLs(seeds)
	if len(allowedDomains) < 1 {
		//nolint:goerr113
		return nil, fmt.Errorf("invalid seed urls")
	}

	// Create the collector settings
	coOptions := []func(*colly.Collector){
		colly.AllowedDomains(allowedDomains...),
		colly.Async(false),
	}

	if o.Verbose {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)

	// Visit each specific folder.
	co.OnHTML(HTMLTagLink, func(e *colly.HTMLElement) {
		href := e.Attr(HTMLAttrRef)

		folderMatch := folderPattern.FindStringSubmatch(href)

		// if the URL is of a folder.
		//nolint:nestif
		if len(folderMatch) > 0 {
			// Do not traverse the hierarchy in reverse order.
			if strings.Contains(href, UpDir) || href == RootDir {
				return
			}

			exactFolderMatch := exactFolderPattern.FindStringSubmatch(href)
			if len(exactFolderMatch) > 0 {
				hrefAbsURL, _ := url.Parse(e.Request.AbsoluteURL(href))

				if !urlSliceContains(seeds, hrefAbsURL) {
					folders = append(folders, path.Base(hrefAbsURL.Path))
					urls = append(urls, hrefAbsURL.String())
				}
			}
			if o.Recursive {
				//nolint:errcheck
				co.Visit(e.Request.AbsoluteURL(href))
			}
		}
	})

	co.OnRequest(func(r *colly.Request) {
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// if the URL is not of a folder.
		if len(folderMatch) == 0 {
			r.Abort()
		}
	})

	// Visit each root folder.
	for _, seedURL := range seeds {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error scraping folder with URL %seedURLs", seedURL.String()))
		}
	}

	return &Result{BaseNames: folders, URLs: urls}, nil
}
