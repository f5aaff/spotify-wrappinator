package requests

import (
	"errors"
	"io/ioutil"
	agent "wrappinator.agent"
)

const (
	baseURL = "https://api.spotify.com/v1/"
)

type ClientRequest struct {
	BaseURL    string
	RequestURL string
	Response   []byte
}

type ClientOpt func(clientRequest *ClientRequest)

func WithBaseURL(url string) ClientOpt {
	return func(request *ClientRequest) {
		request.BaseURL = url
	}
}
func WithRequestURL(url string) ClientOpt {
	return func(request *ClientRequest) {
		request.RequestURL = url
	}
}

func New(clientopts ...ClientOpt) *ClientRequest {
	c := &ClientRequest{}
	for _, opt := range clientopts {
		opt(c)
	}
	return c
}

func GetRequest(a *agent.Agent, request *ClientRequest) {
	get, _ := a.Client.Get(request.BaseURL + request.RequestURL)
	err := errors.New("")
	request.Response, err = ioutil.ReadAll(get.Body)
	if err != nil {
		request.Response = nil
	}
}

func ParamRequest(a *agent.Agent, request *ClientRequest, opts ...RequestOption) {
	fullUrl := request.BaseURL + request.RequestURL
	if params := ProcessOptions(opts...).urlParams.Encode(); params != "" {
		fullUrl += "?" + params
	}

	res, _ := a.Client.Get(fullUrl)
	err := errors.New("")
	request.Response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		request.Response = nil
	}
}
