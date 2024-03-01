package spotify

import (
"context"
"errors"
"net/http"
"encoding/json"
"time"

//"os"

"golang.org/x/oauth2"
)

const (

	rateLimitExceededCode = 429
	defaultRetryDuration = time.Second * 5
)

type Client struct{
	http *http.Client
	baseURL string
	autoRetry bool
	acceptLanguage string
}

func New(httpClient *http.Client,opts ...ClientOption) *Client{
	c := &Client{
		http: httpClient,
		baseURL: "https://api.spotify.com/v1/",
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
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

func retryDuration(resp *http.Response) time.Duration {
	raw := res.Header.Get("Retry-After")
	if raw == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(raw,10,32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second

}




func (c *Client) get(ctx context.Context, url string, result interface{}) error {

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if c.acceptLanguage != "" {
			req.Header.Set("Accept-Language",c.acceptLanguage)
		}
		if err != nil {
			return err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		
		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededCode {
			select {
			case <- ctx.Done():
			case <- time.After(retryDuration(resp)):
				continue
			}
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			return json.NewDecoder(resp).Decode(result)
		}
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}