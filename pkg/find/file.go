package find

import (
	"fmt"
	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
	"github.com/pkg/errors"
	"net/url"
	"path"
	"regexp"
	"strings"
)

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
		colly.Async(true),
	}

	if o.Verbose {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)

	// Add the callback to Visit the linked resource, for each HTML element found
	co.OnHTML(HTMLTagLink, func(e *colly.HTMLElement) {
		href := e.Attr(HTMLAttrRef)

		folderMatch := folderPattern.FindStringSubmatch(href)

		u, _ := url.JoinPath(e.Request.URL.String(), href)

		// If the URL is not of a folder.
		if len(folderMatch) == 0 {
			fileMatch := filePattern.FindStringSubmatch(href)

			// If the URL is of a file.
			if len(fileMatch) > 0 {
				fileName := path.Base(href)
				fileNameMatch := exactFilePattern.FindStringSubmatch(fileName)

				// If the URL matches the file filter regex.
				if len(fileNameMatch) > 0 {
					files = append(files, fileName)
					urls = append(urls, u)
				}
			}
		}

		// Do not traverse the hierarchy in reverse order.
		if o.Recursive && !(strings.Contains(href, UpDir)) && href != RootDir {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(href))
		}
	})

	// Visit each root folder.
	for _, seedURL := range seeds {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error scraping file with URL %seedURLs", seedURL.String()))
		}
	}

	// Wait until colly goroutines are finished.
	co.Wait()

	return &Result{BaseNames: files, URLs: urls}, nil
}
