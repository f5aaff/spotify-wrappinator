package agent

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

func listparam(listIn []string, paramKey string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		values := fmt.Sprintf(strings.Join(listIn[:], ","))
		values = url.QueryEscape(values)
		ro.UrlParams.Set(paramKey, values)
	}
}

func percentParam(inputVal int, paramKey string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		input := inputVal / 100
		if input > 1 {
			return
		}
		ro.UrlParams.Set(paramKey, strconv.Itoa(input))
	}
}

func SeedList(seedList []string, param string) {
	acceptedSeeds := []string{"seed_artists", "seed_genres", "seed_tracks"}
	for _, x := range acceptedSeeds {
		if param == x {
			listparam(seedList, param)
		}
	}
}
func PercentParam(percentage int, paramKey string) {
	acceptedKeys := []string{}
	for _, x := range acceptedKeys {
		if paramKey == x {
			percentParam(percentage, paramKey)
		}
	}
}
