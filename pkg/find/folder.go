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
		colly.Async(o.Async),
		colly.MaxBodySize(o.MaxBodySize),
	}

	if o.Verbose {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)
	if o.ClientTransport != nil {
		co.WithTransport(o.ClientTransport)
	}

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

	// Wait until colly goroutines are finished.
	co.Wait()

	return &Result{BaseNames: folders, URLs: urls}, nil
}
