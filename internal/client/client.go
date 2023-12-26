package client

import (
	"net/http"
)

type HyperstackClient struct {
	client    *http.Client
	apiToken  string
	apiServer string
}

func NewHyperstackClient(
	apiToken string,
	apiServer string,
) *HyperstackClient {
	return &HyperstackClient{
		client:    http.DefaultClient,
		apiToken:  apiToken,
		apiServer: apiServer,
	}
}
