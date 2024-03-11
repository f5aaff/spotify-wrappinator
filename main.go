package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	//random state string, probably has some actual use - but I'm not using it
	state = "abc123"
	//auth is a wrapper of sorts around an oauth2 authenticator struct, allows for context based tokens, so they can be backgrounded and re-authenticate when required.
	a          = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate))
	validToken oauth2.Token
)

func main() {

	//if a token cannot be read from file, request a new token, ask user to login through spotify Oauth
	if readTokenFromFile(&validToken) == nil {
		//set env's are here as a sort of placeholder, needs to be parted out into config file to be filled by user
		err := os.Setenv("SPOTIFY_ID", clientId)
		if err != nil {
			return
		}
		err = os.Setenv("SPOTIFY_SECRET", clientSecret)
		if err != nil {
			return
		}

		http.HandleFunc("/callback", completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("request: ", r.URL.String())
		})

		authURL := a.AuthURL(state)
		fmt.Println("login at this authURL:", authURL)
	}

	//provided a valid token from file, listen on port 8080 for requests to hand off to spotify.
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
}
func StoreTokenToFile(tok *oauth2.Token) {
	f, _ := json.MarshalIndent(tok, "", " ")
	path := tokenStorePath + tokenStoreFileName
	err := ioutil.WriteFile(path, f, 0644)
	if err != nil {
		fmt.Println("broken", err)
		return
	}
}
func readTokenFromFile(tok *oauth2.Token) *oauth2.Token {
	fmt.Println("reading token from file...")
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
func completeAuth(w http.ResponseWriter, r *http.Request) {

	tok, err := a.Token(r.Context(), state, r)
	fmt.Println("storing token to file...")
	StoreTokenToFile(tok)
	if err != nil {
		http.Error(w, "could not retrieve token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	_, err = fmt.Fprintf(w, "login completed! %s", tok)
	if err != nil {
		return
	}
}
func refreshToken(w http.ResponseWriter, r *http.Request) {
	tok, err := a.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "token could not be refreshed", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("state mismatch: %s != %s\n", st, state)
	}

	_, err = fmt.Fprintf(w, "login complete %s", tok)
	if err != nil {
		return
	}
}
