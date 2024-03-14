package agent

import (
	"encoding/json"
	"errors"
	"fmt"
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
	conf  *oauth2.Config
	token *oauth2.Token
}
type AgentOpt func(a *Agent)

func WithToken(token oauth2.Token) AgentOpt {
	return func(a *Agent) {
		a.token = &token
	}
}
func WithConf(conf oauth2.Config) AgentOpt {
	return func(a *Agent) {
		a.conf = &conf
	}
}

func New(conf oauth2.Config, token oauth2.Token, agentopts ...AgentOpt) *Agent {
	a := &Agent{
		conf:  &conf,
		token: &token,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func readTokenFromFile(a *Agent) bool {
	log.Println("readTokenFromFile: reading token from file...")
	f, err := ioutil.ReadFile(tokenStorePath + tokenStoreFileName)
	if err == nil {
		err := json.Unmarshal(f, &a.token)
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

func AuthoriseSession(w http.ResponseWriter, r *http.Request) (*oauth2.Token, error) {
	token, err := auth.GetToken(conf, r.Context(), state, r)
	if err != nil {
		http.Error(w, "token could not be retrieved", http.StatusForbidden)
		return nil, errors.New(err.Error())
	}
	log.Println("AuthoriseSession: Storing Token to File...")
	if err = StoreTokenToFile(token); err != nil {
		log.Println("Could Not Save token:" + err.Error())
		return nil, errors.New(err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		return nil, errors.New("AuthoriseSession: state Mismatch")
	}

	_, err = fmt.Fprintf(w, "login successful\n%s", token)
	if err != nil {
		log.Printf("AuthoriseSession: " + err.Error())
		return nil, errors.New(err.Error())
	}
	return token, nil
}
