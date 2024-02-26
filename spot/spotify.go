package spot

import (
//"context"
"errors"
"net/http"
//"os"

"golang.org/x/oauth2"
)

type Client struct{
	http *http.Client
	baseURL string
	autoRetry bool
	acceptLanguage string
}

type ClientOption func(client *Client)

func (c *Client) Token() (*oauth2.Token, error) {
	transport, ok := c.http.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("spotify: client not backed by oauth2 transport")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}
	return t,nil
}