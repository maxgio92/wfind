package find_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/maxgio92/wfind/internal/network"
	"github.com/maxgio92/wfind/pkg/find"
)

const (
	seedURL         = "https://mirrors.edge.kernel.org/centos/8-stream"
	fileRegexp      = "repomd.xml$"
	expectedResults = 155
)

var _ = Describe("File crawling", func() {
	Context("Async", func() {
		var (
			search = find.NewFind(
				find.WithAsync(true),
				find.WithSeedURLs([]string{seedURL}),
				find.WithClientTransport(network.DefaultClientTransport),
				find.WithFilenameRegexp(fileRegexp),
				find.WithFileType(find.FileTypeReg),
				find.WithRecursive(true),
				find.WithMaxBodySize(1024*512),
				find.WithConnTimeoutRetryBackOff(find.DefaultExponentialBackOffOptions),
				find.WithConnResetRetryBackOff(find.DefaultExponentialBackOffOptions),
			)
			actual        *find.Result
			err           error
			expectedCount = expectedResults
		)
		BeforeEach(func() {
			actual, err = search.Find()
		})
		It("Should not fail", func() {
			Expect(err).To(BeNil())
		})
		It("Should stage results", func() {
			Expect(actual.URLs).ToNot(BeEmpty())
			Expect(actual.URLs).ToNot(BeNil())
		})
		It("Should stage exact result count", func() {
			Expect(len(actual.URLs)).To(Equal(expectedCount))
		})
	})
})
