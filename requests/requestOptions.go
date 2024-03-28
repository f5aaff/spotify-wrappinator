package requests

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
func Fields(fields string) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.UrlParams.Set("fields", fields)
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
