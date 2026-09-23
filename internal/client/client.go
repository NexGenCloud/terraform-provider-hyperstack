package client

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type HyperstackClient struct {
	Client    *http.Client
	ApiToken  string
	ApiServer string

	clientID  string
	userAgent string
}

func NewHyperstackClient(
	apiToken string,
	apiServer string,
	version string,
) *HyperstackClient {
	if version == "" {
		version = "dev"
	}
	// Same shape as hyperstack-sdk-go's own headers. amd64 is spelled x86_64
	// there, so map it here too or the two disagree on Linux.
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	clientID := "terraform-provider-hyperstack/" + version
	return &HyperstackClient{
		Client:    http.DefaultClient,
		ApiToken:  apiToken,
		ApiServer: apiServer,
		clientID:  clientID,
		userAgent: fmt.Sprintf("%s (Go/%s; %s/%s)",
			clientID, strings.TrimPrefix(runtime.Version(), "go"), runtime.GOOS, arch),
	}
}

func (c HyperstackClient) GetAddHeadersFn() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("api_key", c.ApiToken)
		// TODO: do we need to support it?
		//req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}

// identityTransport stamps the provider's identity onto every outgoing request.
type identityTransport struct {
	base      http.RoundTripper
	clientID  string
	userAgent string
}

func (t identityTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Hyperstack-Client", t.clientID)
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}

// InstallIdentityTransport wraps http.DefaultTransport so every SDK client
// picks it up.
func (c HyperstackClient) InstallIdentityTransport() {
	if t, ok := http.DefaultTransport.(identityTransport); ok {
		t.clientID = c.clientID
		t.userAgent = c.userAgent
		http.DefaultTransport = t
		return
	}
	http.DefaultTransport = identityTransport{
		base:      http.DefaultTransport,
		clientID:  c.clientID,
		userAgent: c.userAgent,
	}
}
