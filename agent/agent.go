package agent

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	auth "wrappinator.auth"
)

const (
	redirectURL         = "http://localhost:8080/callback"
	clientId     string = "1b0ac2b304e941d9890dc016171c2226"
	clientSecret string = "dd8f644ef4074f7f82daca80487818b6"
	//this needs templating out badly, can't be having hard coded directories in this
	tokenStorePath     string = "/home/f5adff/.config/wrappinator/token/"
	tokenStoreFileName string = "token.json"
)

var (
	state      string = "abc123"
	conf              = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate))
	validToken oauth2.Token
)

type Agent struct {
	Conf   *oauth2.Config
	Token  *oauth2.Token
	Client *http.Client
}
type AgentOpt func(a *Agent)
type AgentFunc func(a *Agent)

func WithToken(token oauth2.Token) AgentOpt {
	return func(a *Agent) {
		a.Token = &token
	}
}
func WithConf(conf oauth2.Config) AgentOpt {
	return func(a *Agent) {
		a.Conf = &conf
	}
}

func WithClient(client *http.Client) AgentOpt {
	return func(a *Agent) {
		a.Client = client
	}
}

func GetClient() AgentOpt {
	return func(a *Agent) {
		if a.Client != nil && a.Conf != nil && a.Token != nil {
			a.Client = auth.Client(a.Conf, context.Background(), a.Token)
		}
	}
}

func New(conf *oauth2.Config, agentOpts ...AgentOpt) *Agent {
	a := &Agent{
		Conf: conf,
	}
	for _, opt := range agentOpts {
		opt(a)
	}
	return a
}

func ReadTokenFromFile(a *Agent) bool {
	log.Println("readTokenFromFile: reading token from file...")
	f, err := ioutil.ReadFile(tokenStorePath + tokenStoreFileName)
	if err == nil {
		err := json.Unmarshal(f, &a.Token)
		if err != nil {
			log.Println("readTokenFromFile: ERROR: " + err.Error())
			return false
		}
		log.Println("token read from file successfully...")
	}
	return true
}
func StoreTokenToFile(tok *oauth2.Token) error {
	f, _ := json.MarshalIndent(tok, "", " ")
	path := tokenStorePath + tokenStoreFileName
	err := ioutil.WriteFile(path, f, 0644)
	if err != nil {

		log.Println("StoreTokenToFile error:", err)
		return errors.New("StoreTokenToFile:" + err.Error())
	}
	return nil
}
