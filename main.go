package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	agent "wrappinator.agent"
	auth "wrappinator.auth"
	device "wrappinator.device"
	requests "wrappinator.requests"
)

const (
	redirectURL = "http://localhost:8080/callback"
)

var (
	envloaderr          = godotenv.Load()
	state        string = "abc123"
	clientId     string = os.Getenv("CLIENT_ID")
	clientSecret string = os.Getenv("CLIENT_SECRET")
	conf                = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate, auth.ScopeUserReadPlaybackState, auth.ScopeUserModifyPlaybackState, auth.ScopeStreaming))
	validToken   oauth2.Token
	a            = agent.New(conf, agent.WithToken(validToken))
	d            = device.New()
)

func main() {

	if envloaderr != nil {
		return
	}
	/*
		if a token can't be read from file, prompt the user to log in
	*/
	if agent.ReadTokenFromFile(a) == false {

		http.HandleFunc("/callback", AuthoriseSession)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("request: ", r.URL.String())
		})

		url := auth.GetAuthURL(conf, state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := errors.New("")
	a.Token, err = auth.RefreshToken(conf, context.Background(), a.Token)
	if err != nil {
		log.Fatalf("Token Refresh error: %s", err.Error())
	}
	a.Client = auth.Client(conf, context.Background(), a.Token)

	err = d.GetCurrentDevice(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", d)

	c := requests.New()
	err = d.ChangeVolume(a, c, 10)
	if err != nil {
		fmt.Println("broke:" + err.Error())
		return
	}
	fmt.Printf("%+v\n", c)
	err = d.GetCurrentDevice(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", d)

	fmt.Printf("%+v\n", a.Client)
}

func AuthoriseSession(w http.ResponseWriter, r *http.Request) {
	err := errors.New("")
	a.Token, err = auth.GetToken(a.Conf, r.Context(), state, r)
	if err != nil {
		http.Error(w, "token could not be retrieved", http.StatusForbidden)
		log.Fatal(err)
	}
	log.Println("AuthoriseSession: Storing Token to File...")
	if err = agent.StoreTokenToFile(a.Token); err != nil {
		log.Println("Could Not Save token:" + err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("state mismatch: %s != %s\n", st, state)
	}

	_, err = fmt.Fprintf(w, "login successful\n%s", a.Token)
	if err != nil {
		log.Printf("AuthoriseSession: " + err.Error())
		return
	}
}
