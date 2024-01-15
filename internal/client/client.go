package client

import (
	"context"
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

func (c HyperstackClient) GetAddHeadersFn() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("api_key", c.ApiToken)
		// TODO: do we need to support it?
		//req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}
