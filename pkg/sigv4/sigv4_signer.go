package sigv4

/*
Copyright 2022 MOIA GmbH

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

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

type SigV4Transport struct {
	signer *v4.Signer
	config *Config
	next   http.RoundTripper
}

type Config struct {
	Service string
	Region  string
}

// The RoundTripperFunc type is an adapter to allow the use of ordinary
// functions as RoundTrippers. If f is a function with the appropriate
// signature, RoundTripperFunc(f) is a RoundTripper that calls f.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the RoundTripper interface.
func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

// NewSigner instantiates a new signing middleware with an optional succeeding
// middleware. The http.DefaultTransport will be used if nil.
func NewSigner(cfg *Config, creds *credentials.Credentials, next http.RoundTripper) (http.RoundTripper, error) {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if next == nil {
			next = http.DefaultTransport
		}

		signer, err := createSigner(creds, false)
		if err != nil {
			return nil, err
		}
		m := &SigV4Transport{
			config: cfg,
			next:   next,
			signer: signer,
		}
		return m.exec(r)
	}), nil
}

func (m *SigV4Transport) exec(origReq *http.Request) (*http.Response, error) {
	req, err := m.createSignedRequest(origReq)
	if err != nil {
		return nil, err
	}

	//nolint: wrapcheck
	return m.next.RoundTrip(req)
}

func (m *SigV4Transport) createSignedRequest(origReq *http.Request) (*http.Request, error) {
	req, err := http.NewRequest(origReq.Method, origReq.URL.String(), origReq.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP req: %w", err)
	}

	body := bytes.NewReader([]byte{})
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read body to sign: %w", err)
		}
		body = bytes.NewReader(b)
	}

	if strings.Contains(req.URL.RawPath, "%2C") {
		req.URL.RawPath = rest.EscapePath(req.URL.RawPath, false)
	}

	_, err = m.signer.Sign(req, body, m.config.Service, m.config.Region, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("signing body failed: %w", err)
	}

	copyHeaderWithoutOverwrite(req.Header, origReq.Header)

	return req, nil
}

func createSigner(c *credentials.Credentials, verboseMode bool) (*v4.Signer, error) {
	signerOpts := func(s *v4.Signer) {
		if verboseMode {
			s.Debug = aws.LogDebugWithSigning
		}
	}

	return v4.NewSigner(c, signerOpts), nil
}

func copyHeaderWithoutOverwrite(dst, src http.Header) {
	for k, vv := range src {
		if _, ok := dst[k]; !ok {
			for _, v := range vv {
				dst.Add(k, v)
			}
		}
	}
}
