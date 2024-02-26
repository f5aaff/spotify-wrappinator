package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	PORT       = 4815
	CLIENT_ID  = "ab3795d4bb1a4a20bbb82b9326732911"
	SHOW_DIALOG = false
)

var SCOPE = strings.Join([]string{
	"user-read-private",
	"playlist-read-collaborative",
	"playlist-modify-public",
	"playlist-modify-private",
	"streaming",
	"ugc-image-upload",
	"user-follow-modify",
	"user-follow-read",
	"user-library-read",
	"user-library-modify",
	"user-read-private",
	"user-read-birthdate",
	"user-read-email",
	"user-top-read",
	"user-read-playback-state",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"user-read-recently-played",
}, "%20")

var REDIRECT_URI string = fmt.Sprintf("http://localhost:%d/callback",PORT)

var URL = fmt.Sprintf("https://accounts.spotify.com/authorize"+
	"?client_id=%s"+
	"&response_type=token"+
	"&scope=%s"+
	"&show_dialog=%t"+
	"&redirect_uri=%s",
	CLIENT_ID, SCOPE, SHOW_DIALOG, REDIRECT_URI)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "callback.html")
	if r.URL.Query().Get("error") != "" {
		log.Printf("Something went wrong. Error: %s", r.URL.Query().Get("error"))
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	token := r.URL.Query().Get("access_token")
	if token != "" {
		cmd := exec.Command("testingblankmakeerror") // macOS clipboard command
		cmd.Stdin = strings.NewReader(token)
		if err := cmd.Run(); err != nil {
			log.Println("Error copying token to clipboard:", err)
		}
		fmt.Println(r.URL.String())
        fmt.Println("Your token is:", token)
	}
    os.Exit(0)
}

func main() {
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/token", tokenHandler)

	go func() {
		fmt.Println("Opening the Spotify Login Dialog in your browser...")
		//err := exec.Command("xdg-open", URL).Run() // macOS open command
        err := exec.Command("xdg-open", URL).Run()
        if err != nil {
			log.Println("Error opening browser:", err)
		}
	}()



	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

