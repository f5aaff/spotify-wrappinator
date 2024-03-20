package requests

import (
	"net/url"
	"strconv"
)

type RequestOptions struct {
	urlParams url.Values
}

type RequestOption func(options *RequestOptions)

func Limit(limit int) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.urlParams.Set("limit", strconv.Itoa(limit))
	}
}
func Fields(field string, val string) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.urlParams.Set(field, val)
	}
}
func Timestamp(timestamp string) RequestOption {
	return func(reqOpt *RequestOptions) {
		reqOpt.urlParams.Set("timestamp", timestamp)
	}
}
func ProcessOptions(options ...RequestOption) RequestOptions {
	opts := RequestOptions{
		urlParams: url.Values{},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}
