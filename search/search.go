package agent

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	auth "wrappinator.auth"
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


