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
				find.WithMaxBodySize(find.DefaultMaxBodySize),
				find.WithContextDeadlineRetryBackOff(find.DefaultExponentialBackOffOptions),
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
