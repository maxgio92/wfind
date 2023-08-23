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

import "net/url"

func getHostnamesFromURLs(urls []*url.URL) []string {
	hostnames := []string{}

	for _, v := range urls {
		hostnames = append(hostnames, v.Host)
	}

	return hostnames
}

func urlSliceContains(us []*url.URL, u *url.URL) bool {
	for _, v := range us {
		if v == u {
			return true
		}
	}

	return false
}
