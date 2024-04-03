package wrappinator.requests

import (
	"net/url"
	"strconv"
)

type RequestOptions struct {
	UrlParams url.Values
}

type RequestOption func(options *RequestOptions)

func Limit(limit int) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.UrlParams.Set("limit", strconv.Itoa(limit))
	}
}
func Fields(field string, value string) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.UrlParams.Set(field, value)
	}
}
func Timestamp(timestamp string) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.UrlParams.Set("timestamp", timestamp)
	}
}
func ProcessOptions(options ...RequestOption) RequestOptions {
	opts := RequestOptions{
		UrlParams: url.Values{},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}
