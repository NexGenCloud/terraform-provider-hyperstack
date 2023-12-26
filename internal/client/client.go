package client

import (
	"net/http"
)

type HyperstackClient struct {
	Client    *http.Client
	ApiToken  string
	ApiServer string
}

func NewHyperstackClient(
	apiToken string,
	apiServer string,
) *HyperstackClient {
	return &HyperstackClient{
		Client:    http.DefaultClient,
		ApiToken:  apiToken,
		ApiServer: apiServer,
	}
}
