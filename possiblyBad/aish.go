
 package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

const (
    spotifyAPIBaseURL = "https://api.spotify.com/v1"
    clientID          = "1b0ac2b304e941d9890dc016171c2226"
    clientSecret      = "dd8f644ef4074f7f82daca80487818b6"
    redirectURI       = "localhost:8080"
)

func main() {
    // Get an access token
    token, err := getAccessToken()
    if err != nil {
        fmt.Println("Error getting access token:", err)
        return
    }
    // Example: Get a user's profile
    userProfile, err := getUserProfile(token)
    if err != nil {
        fmt.Println("Error fetching user profile:", err)
        return
    }
    fmt.Println(userProfile)
    fmt.Printf("User ID: %s\nDisplay Name: %s\n", userProfile.ID, userProfile.DisplayName)

    playlists,err:= getUserPlaylists(token)
    if err != nil{
        fmt.Println("Error fetching playlists:", err)
        return
    }

    for playlist := range playlists{
        fmt.Println(playlist)
    }
}




func getAccessToken() (string, error) {
    // Base64 encode the client ID and secret
    authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

    // Create a POST request to get the access token
    reqBody := strings.NewReader("grant_type=client_credentials")
    req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", reqBody)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Basic "+authHeader)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Parse the JSON response
    // (You should use a proper JSON library for production code)
    accessToken := string(body)
    return accessToken, nil
}

func getUserProfile(token string) (*UserProfile, error) {
    // Create a GET request to fetch the user profile
    req, err := http.NewRequest("GET", spotifyAPIBaseURL+"/me", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+token)

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse the JSON response
    // (Again, use a proper JSON library in production)
    var userProfile UserProfile
    err = json.NewDecoder(resp.Body).Decode(&userProfile)
    if err != nil {
        return nil, err
    }

    return &userProfile, nil
}





// Define a struct to hold user profile data
type UserProfile struct {
    ID           string `json:"id"`
    DisplayName  string `json:"display_name"`
    // Add other relevant fields as needed
}

func getUserPlaylists(token string) ([]Playlist, error) {
    // Create a GET request to fetch user playlists
    req, err := http.NewRequest("GET", spotifyAPIBaseURL+"/me/playlists", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+token)

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse the JSON response
    var playlistsResponse struct {
        Items []Playlist `json:"items"`
    }
    err = json.NewDecoder(resp.Body).Decode(&playlistsResponse)
    if err != nil {
        return nil, err
    }

    return playlistsResponse.Items, nil
}
func getRefreshToken(authCode string) (string, error) {
    // Base64 encode the client ID and secret
    authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

    // Create a POST request to exchange the authorization code for tokens
    reqBody := strings.NewReader(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", authCode, redirectURI))
    req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", reqBody)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Basic "+authHeader)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Parse the JSON response
    var tokenResponse struct {
        RefreshToken string `json:"refresh_token"`
    }
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        return "", err
    }

    return tokenResponse.RefreshToken, nil
}
// Define a struct to hold playlist data
type Playlist struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    // Add other relevant fields as needed
}

