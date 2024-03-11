package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	auth "wrappinator.auth"
)

const (
	redirectURL               = "http://localhost:8080/callback"
	clientId           string = "1b0ac2b304e941d9890dc016171c2226"
	clientSecret       string = "dd8f644ef4074f7f82daca80487818b6"
	tokenStorePath     string = "/home/f5adff/.config/wrappinator/token/"
	tokenStoreFileName string = "token.json"
)

var (
	state      = "abc123"
	a          = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate))
	validToken oauth2.Token
)

type getRequest struct {
	authorisation  string
	targetEndPoint string
	variable       string
	value          any
}

type fetchRequest struct {
	authorisation  string
	targetEndPoint string
	body           string
}

type postRequest struct {
	authorisation  string
	targetEndPoint string
	body           url.Values
}

func main() {

	err := os.Setenv("SPOTIFY_ID", clientId)
	if err != nil {
		return
	}
	err = os.Setenv("SPOTIFY_SECRET", clientSecret)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", os.Getenv("SPOTIFY_ID"))
	fmt.Printf("%s\n", os.Getenv("SPOTIFY_SECRET"))

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: ", r.URL.String())
	})

	authURL := a.AuthURL(state)
	fmt.Println("login at this authURL:", authURL)

	err = http.ListenAndServe(":8080", nil)
	fmt.Println("aaaaaaaaaa", validToken)
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
	f, err := ioutil.ReadFile(tokenStorePath + tokenStoreFileName)
	if err == nil {
		err := json.Unmarshal(f, &tok)
		if err != nil {
			fmt.Println("readTokenFromFile:", err)
			return nil
		}
		return tok
	}
	return nil
}
func completeAuth(w http.ResponseWriter, r *http.Request) {

	if readTokenFromFile(&validToken) == nil {
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
func sendGetRequest(req *getRequest, client http.Client) *http.Response {
	endPoint := req.targetEndPoint
	body := fmt.Sprintf("%s/%s", req.variable, req.value)
	call := fmt.Sprintf("%s/%s", endPoint, body)

	request, err := http.NewRequest("GET", call, nil)
	request.Header = http.Header{
		"Authorisation": {req.authorisation},
	}
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return res
	}
	return nil
}

func sendPostRequest(req *postRequest, client http.Client) *http.Response {

	r, _ := http.NewRequest("POST", req.targetEndPoint, bytes.NewBufferString(req.body.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	//fmt.Println(r)
	return resp
}

func getToken(clientId string, clientSecret string, client http.Client) string {
	req := postRequest{
		authorisation:  "",
		targetEndPoint: "https://accounts.spotify.com/api/token",
		body: url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {clientId},
			"client_secret": {clientSecret},
		},
	}

	resp := sendPostRequest(&req, client)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, mapErr := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	if mapErr != nil {
		fmt.Println("error mapping string", mapErr)
	}
	jsonErr := json.Unmarshal(body, &res)
	if jsonErr != nil {
		fmt.Println("Error Parsing JSON:", jsonErr)
	}
	accessToken, ok := res["access_token"].(string)
	if !ok {
		fmt.Println("access token not found in response")
	}
	return accessToken
}
