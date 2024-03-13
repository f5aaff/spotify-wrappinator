package test

import (
	auth "authNoStruct"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {

	/*
		if a token can't be read from file, prompt the user to log in
	*/
	if readTokenFromFile(&validToken) == nil {

		http.HandleFunc("/callback", AuthoriseSession)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("request: ", r.URL.String())
		})
	}
	token, err := auth.RefreshToken(conf, context.Background(), &validToken)
	if err != nil {
		log.Fatalf("Token Refresh error: %s", err.Error())
	}
	fmt.Println(token)
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
func readTokenFromFile(tok *oauth2.Token) *oauth2.Token {
	log.Println("readTokenFromFile: reading token from file...")
	f, err := ioutil.ReadFile(tokenStorePath + tokenStoreFileName)
	if err == nil {
		err := json.Unmarshal(f, &tok)
		if err != nil {
			return nil
		}
		return tok
	}
	return nil
}

func AuthoriseSession(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(conf, r.Context(), state, r)
	if err != nil {
		http.Error(w, "token could not be retrieved", http.StatusForbidden)
		log.Fatal(err)
	}
	log.Println("AuthoriseSession: Storing Token to File...")
	if err = StoreTokenToFile(token); err != nil {
		log.Println("Could Not Save token:" + err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("state mismatch: %s != %s\n", st, state)
	}

	_, err = fmt.Fprintf(w, "login successful\n%s", token)
	if err != nil {
		log.Printf("AuthoriseSession: " + err.Error())
		return
	}
}
