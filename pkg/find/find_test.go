package find_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/reply"

	"github.com/maxgio92/wfind/pkg/find"
)

const (
	homedir  string = "home"
	filename string = "File"
	dirname  string = "Dir"
)

var (
	files       = []string{"hello", "world"}
	subdirs     = []string{"foo", "bar", "baz"}
	homedirBody = fmt.Sprintf(`
<html>
<head><title>Index of %s/</title></head>
<body>
<h1>Index of /%s/</h1><hr><pre><a href="../">../</a>
<a href="%s/">%s/</a>
<a href="%s/">%s/</a>
<a href="%s/">%s/</a>
<a href="%s">%s/</a>
<a href="%s">%s/</a>
</pre><hr></body>
</html>
`, homedir,
		homedir,
		subdirs[0], subdirs[0],
		subdirs[1], subdirs[1],
		subdirs[2], subdirs[2],
		files[0], files[0],
		files[1], files[1])
	subdirBodyF = `
<html>
<head><title>Index of %s/%s/</title></head>
<body>
<h1>Index of /%s/%s/</h1><hr><pre><a href="../">../</a>
<a href="%s/">%s/</a>
<a href="%s">%s</a>
</pre><hr></body>
</html>`
	subdirBodyDotSlashF = `
<html>
<head><title>Index of %s/%s/</title></head>
<body>
<h1>Index of /%s/%s/</h1><hr><pre><a href="../">../</a>
<a href="./%s/">%s/</a>
<a href="./%s">%s</a>
</pre><hr></body>
</html>`
	subdirBodies         []string
	subdirBodiesDotSlash []string
)

func initFileHierarchy() {
	for _, v := range subdirs {
		subdirBodies = append(subdirBodies,
			fmt.Sprintf(subdirBodyF, homedir, v, homedir, v, dirname, dirname, filename, filename),
		)

		subdirBodiesDotSlash = append(subdirBodiesDotSlash,
			fmt.Sprintf(subdirBodyDotSlashF, homedir, v, homedir, v, dirname, dirname, filename, filename),
		)
	}
}

func initWebServer(t *testing.T) *mocha.Mocha {
	m := mocha.New(t).CloseOnCleanup(t)
	m.Start()

	m.AddMocks(
		// home dir.
		mocha.Get(expect.URLPath(fmt.Sprintf("/%s/", homedir)).
			Or(expect.URLPath(fmt.Sprintf("/%s", homedir)))).
			Reply(reply.OK().BodyString(homedirBody)))

	// Sub directories.
	for i := range subdirs {
		m.AddMocks(
			// File in root.
			mocha.Get(expect.URLPath(fmt.Sprintf("/%s/%s/", homedir, subdirs[i])).
				Or(expect.URLPath(fmt.Sprintf("/%s/%s", homedir, subdirs[i])))).
				Reply(reply.OK().BodyString(subdirBodies[0])),
			// Sub directory.
			mocha.Get(expect.URLPath(fmt.Sprintf("/%s/%s/", homedir, subdirs[i])).
				Or(expect.URLPath(fmt.Sprintf("/%s/%s", homedir, subdirs[i])))).
				Reply(reply.OK().BodyString(subdirBodies[0])),
			// File in sub directory.
			mocha.Get(expect.URLPath(fmt.Sprintf("/%s/%s/%s/", homedir, subdirs[i], dirname)).
				Or(expect.URLPath(fmt.Sprintf("/%s/%s/%s", homedir, subdirs[i], dirname)))).
				Reply(reply.OK().BodyString("")),
			// Directory in sub directory.
			mocha.Get(expect.URLPath(fmt.Sprintf("/%s/%s/%s", homedir, subdirs[i], filename))).
				Reply(reply.OK().BodyString("")))
	}

	return m
}

//nolint:dupl
func TestFindFile(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(`.+`),
		find.WithFileType(find.FileTypeReg),
		find.WithRecursive(false),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(files))
	assert.Equal(t, found.BaseNames, files)
}

//nolint:dupl
func TestFindDir(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(`^.+$`),
		find.WithFileType(find.FileTypeDir),
		find.WithRecursive(false),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(subdirs))
	assert.Equal(t, found.BaseNames, subdirs)
}

//nolint:dupl
func TestFindFileRecursive(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(fmt.Sprintf("^%s$", filename)),
		find.WithFileType(find.FileTypeReg),
		find.WithRecursive(true),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(subdirs))
}

//nolint:dupl
func TestFindDirRecursive(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(fmt.Sprintf("^%s$", dirname)),
		find.WithFileType(find.FileTypeDir),
		find.WithRecursive(true),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(subdirs))
}

//nolint:dupl
func TestFindFileRecursiveDotSlash(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(fmt.Sprintf("^%s$", filename)),
		find.WithFileType(find.FileTypeReg),
		find.WithRecursive(true),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(subdirs))
}

//nolint:dupl
func TestFindDirRecursiveDotSlash(t *testing.T) {
	t.Parallel()

	initFileHierarchy()
	m := initWebServer(t)

	finder := find.NewFind(
		find.WithSeedURLs([]string{fmt.Sprintf("%s/%s", m.URL(), homedir)}),
		find.WithFilenameRegexp(fmt.Sprintf("^%s$", dirname)),
		find.WithFileType(find.FileTypeDir),
		find.WithRecursive(true),
		find.WithVerbosity(false),
	)

	found, err := finder.Find()

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.URLs, len(subdirs))
}
