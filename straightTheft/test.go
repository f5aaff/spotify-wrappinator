package main

import (
	"fmt"
	"log"
	"net/http"
//	"context"
//	"encoding/json"
	

	spotifyauth "wrappinator.spotifyauth"
	spotify "wrappinator.spotify"
)

const (
	redirectURI = "localhost:8080/callback"

	)
var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: ", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	} ()

	url := auth.AuthURL(state)
	fmt.Println("login at: ",url)

	//	targetURL := client.baseURL + "me"
//	var obj map[string]*json.RawMessage
	
	user,err := spotify.get(context.Background(),url,&obj)
//	fmt.Println(user)
//	if err != nil {
//		log.Fatal(err)
//	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
