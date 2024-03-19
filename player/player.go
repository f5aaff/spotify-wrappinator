package player

import (
	"golang.org/x/oauth2"
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
