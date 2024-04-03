package spotify_wrappinator.search

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"strconv"
	"strings"
	auth "wrappinator.auth"
	requests "wrappinator.requests"
)

const (
	redirectURL = "http://localhost:8080/callback"
)

var (
	envloaderr                = godotenv.Load()
	tokenStorePath     string = os.Getenv("TOKEN_STORE_PATH")
	tokenStoreFileName string = os.Getenv("TOKEN_FILENAME")
	clientId           string = os.Getenv("CLIENT_ID")
	clientSecret       string = os.Getenv("CLIENT_SECRET")
	state              string = "abc123"
	conf                      = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate))
	validToken         oauth2.Token
)

func Query(query string, tagvals map[string]string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		acceptedTags := []string{"album", "artist", "track", "year", "upc", "hipster", "new", "isrc", "genre"}
		var tagString string
		for k, v := range tagvals {
			for _, a := range acceptedTags {
				if k == a {
					tag := fmt.Sprintf("%s:%s", k, v)
					tagString = fmt.Sprintf("%s,%s", tagString, tag)
				}
			}
		}
		query = url.QueryEscape(query)
		tagString = url.QueryEscape(tagString)
		res := fmt.Sprintf("%s\\% %s", query, tagString)
		ro.UrlParams.Set("q", res)
	}
}

func Types(types []string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		input := fmt.Sprintf(strings.Join(types[:], ","))
		input = url.QueryEscape(input)
		ro.UrlParams.Set("type", input)
	}
}
func Market(market string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		ro.UrlParams.Set("market", market)
	}
}
func Offset(offset int) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		ro.UrlParams.Set("limit", strconv.Itoa(offset))
	}
}
func IncludeExternal(audio bool) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		if audio {
			ro.UrlParams.Set("include_external", "audio")
		}
	}
}
