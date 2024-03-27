package requests

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	agent "wrappinator.agent"
)

const (
	baseURL = "https://api.spotify.com/v1/"
)

type ClientRequest struct {
	BaseURL    string
	RequestURL string
	Response   []byte
	Method     string
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
	if params := ProcessOptions(opts...).UrlParams.Encode(); params != "" {
		fullUrl += "?" + params
	}

	res, _ := a.Client.Get(fullUrl)
	err := errors.New("")
	request.Response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		request.Response = nil
	}
}

func PutRequest(a *agent.Agent, c *ClientRequest, opts ...RequestOption) {
	fullUrl := c.BaseURL + c.RequestURL
	if params := ProcessOptions(opts...).UrlParams.Encode(); params != "" {
		fullUrl += "?" + params
	}
	requrl, err := url.Parse(fullUrl)
	if err != nil {
		return
	}
	head := http.Header{}
	bearerval := "bearer " + a.Token.AccessToken
	head.Set("Authorization", bearerval)
	req := http.Request{Method: "PUT", URL: requrl, Header: head}
	res, err := a.Client.Do(&req)
	if err != nil {
		c.Response = []byte(err.Error())
		return
	}
	c.Response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		c.Response = []byte(err.Error())
	}
}

func ParamFormRequest(a *agent.Agent, c *ClientRequest, opts ...RequestOption) {
	fullUrl := c.BaseURL + c.RequestURL
	if params := ProcessOptions(opts...).UrlParams; params != nil {
		res, _ := a.Client.PostForm(fullUrl, params)
		err := errors.New("")
		c.Response, err = ioutil.ReadAll(res.Body)
		if err != nil {
			c.Response = nil
		}
	}
}
