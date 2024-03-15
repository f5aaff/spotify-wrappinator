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
type ClientGetFunc func(clientRequest *ClientRequest, a *agent.Agent)

func GetReq() ClientGetFunc {
	return func(request *ClientRequest, a *agent.Agent) {
		get, _ := a.Client.Get(request.BaseURL + request.RequestURL)
		err := errors.New("")
		request.Response, err = ioutil.ReadAll(get.Body)
		if err != nil {
			request.Response = nil
		}
	}
}
