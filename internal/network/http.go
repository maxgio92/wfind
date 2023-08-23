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

package network

import (
	"net"
	"net/http"
	"time"
)

const (
	// DefaultTimeout is the default HTTP client transport's timeout in
	// milliseconds.
	// By default, is set to 180 seconds.
	DefaultTimeout = 180000

	// DefaultKeepAlive is the default interval in milliseconds between
	// keep-alive probes for an active network connection.
	// By defdault, is set to 30 seconds.
	DefaultKeepAlive = 30000

	// DefaultMaxIdleConns is the default maximum number of idle
	// (keep-alive) connections across all hosts.
	// By default, is set to 1000.
	DefaultMaxIdleConns = 1000

	// DefaultMaxIdleConnsPerHost is the default maximum number of idle
	// (keep-alive) connections to keep per-host.
	// By default, is set to 1000.
	DefaultMaxIdleConnsPerHost = 1000

	// DefaultIdleConnTimeout is the default maximum amount of time in milliseconds an
	// idle (keep-alive) connection will remain idle before closing itself.
	// By default, is set to 120 seconds.
	DefaultIdleConnTimeout = 120000

	// DefaultTLSHandshakeTimeout is the default maximum amount of time waiting to
	// wait for a TLS handshake.
	// By default, is set to 30 seconds.
	DefaultTLSHandshakeTimeout = 30000
)

var DefaultClientTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   DefaultTimeout * time.Millisecond,
		KeepAlive: DefaultKeepAlive * time.Millisecond,
	}).DialContext,
	MaxIdleConns:        DefaultMaxIdleConns,
	MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
	IdleConnTimeout:     DefaultIdleConnTimeout * time.Millisecond,
	TLSHandshakeTimeout: DefaultTLSHandshakeTimeout * time.Millisecond,
}

type TransportOption func(o *http.Transport)

func WithDialer(dialer *net.Dialer) TransportOption {
	return func(o *http.Transport) {
		o.DialContext = dialer.DialContext
	}
}

func WithMaxIdleConns(maxIdleConns int) TransportOption {
	return func(o *http.Transport) {
		o.MaxIdleConns = maxIdleConns
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) TransportOption {
	return func(o *http.Transport) {
		o.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func WithIdleConnsTimeout(idleConnsTimeout time.Duration) TransportOption {
	return func(o *http.Transport) {
		o.IdleConnTimeout = idleConnsTimeout
	}
}

func WithTLSHandshakeTimeout(TLSHandshakeTimeout time.Duration) TransportOption {
	return func(o *http.Transport) {
		o.TLSHandshakeTimeout = TLSHandshakeTimeout
	}
}

func NewTransport(opts ...TransportOption) *http.Transport {
	t := &http.Transport{}

	for _, f := range opts {
		f(t)
	}

	return t
}

type DialerOption func(o *net.Dialer)

func WithTimeout(timeout time.Duration) DialerOption {
	return func(o *net.Dialer) {
		o.Timeout = timeout
	}
}

func WithKeepAlive(interval time.Duration) DialerOption {
	return func(o *net.Dialer) {
		o.KeepAlive = interval
	}
}

func NewDialer(opts ...DialerOption) *net.Dialer {
	d := &net.Dialer{}

	for _, f := range opts {
		f(d)
	}

	return d
}
