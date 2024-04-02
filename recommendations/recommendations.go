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
func intParam(inputVal int, paramKey string) requests.RequestOption {
	return func(ro *requests.RequestOptions) {
		ro.UrlParams.Set(paramKey, strconv.Itoa(inputVal))
	}
}
func SeedList(seedList []string) {
	acceptedSeeds := []string{"seed_artists", "seed_genres", "seed_tracks"}
	input := []string{}
	for x, y := range acceptedSeeds {
		if seedList[x] == y {
			input = append(input, seedList[x])
		}
		listparam(input, "seed_list")
	}
}

func ListParams(inputMap map[string][]string) requests.RequestOption {
	acceptedKeys := []string{"seed_artists", "seed_genres", "seed_tracks"}

	for _, key := range acceptedKeys {
		input, ok := inputMap[key]
		if ok {
			listparam(input, key)
		}
	}
	return nil
}
func PercentParams(inputMap map[string]int) requests.RequestOption {
	acceptedKeys := []string{"min_acousticness", "max_acousticness", "target_acousticness", "min_danceability", "max_danceability", "target_danceability", "min_duration_ms", "max_duration_ms", "target_duration_ms", "min_energy", "max_energy", "target_energy", "min_instrumentalness", "max_instrumentalness", "target_instrumentalness", "min_keymax_key", "target_key", "min_liveness", "max_liveness", "target_liveness", "min_loudness", "max_loudness", "target_loudness", "min_mode", "max_mode", "target_mode", "min_popularity", "max_popularity", "target_popularity", "min_speechiness", "max_speechiness", "target_speechiness", "min_valence", "max_valence", "target_valence"}
	for _, y := range acceptedKeys {
		input, ok := inputMap[y]
		if ok {
			percentParam(input, y)
		}
	}
	return nil
}

func IntParams(inputMap map[string]int) requests.RequestOption {
	acceptedKeys := []string{"min_tempo", "max_tempo", "target_tempo", "min_time_signature", "max_time_signature", "target_time_signature"}
	for _, y := range acceptedKeys {
		input, ok := inputMap[y]
		if ok {
			intParam(input, y)
		}
	}
	return nil
}
