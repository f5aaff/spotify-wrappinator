package main

import(
    "fmt"
    "log"
    "net/http"
    "net/url"
    "bytes"
    "encoding/json"
    "io/ioutil"
)

const clientId string = "1b0ac2b304e941d9890dc016171c2226"
const clientSecret string = "dd8f644ef4074f7f82daca80487818b6"

type getRequest struct {
    authorisation string
    targetEndPoint string
    variable string
    value any
}

type fetchRequest struct {
    authorisation string
    targetEndPoint string
    body string

}

type postRequest struct {
    authorisation string
    targetEndPoint string
    body url.Values
}

func sendGetRequest(req *getRequest,client http.Client)(*http.Response){
    endPoint := req.targetEndPoint
    body := fmt.Sprintf("%s/%s",req.variable,req.value)
    call := fmt.Sprintf("%s/%s",endPoint,body)

    request, err := http.NewRequest("GET",call,nil)
    request.Header = http.Header{
        "Authorisation": {req.authorisation},
    }
    res, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
        return nil
    }else{return res}
    return nil
    }



func sendPostRequest(req *postRequest,client http.Client)(*http.Response){

    r,_:=http.NewRequest("POST", req.targetEndPoint, bytes.NewBufferString(req.body.Encode()))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp,_:=client.Do(r)
    //fmt.Println(r)
    return resp
    }

func getToken(clientId string,clientSecret string, client http.Client)(string){
    req := postRequest{
        authorisation : "",
        targetEndPoint: "https://accounts.spotify.com/api/token",
        body : url.Values{
        "grant_type":    {"client_credentials"},
        "client_id":     {clientId},
        "client_secret": {clientSecret},
        },
    }

    resp := sendPostRequest(&req,client)
    defer resp.Body.Close()

    body,mapErr := ioutil.ReadAll(resp.Body)
    var res map[string]interface{}
    if mapErr != nil{
        fmt.Println("error mapping string",mapErr)
    }
    jsonErr := json.Unmarshal(body, &res)
    if jsonErr != nil{
        fmt.Println("Error Parsing JSON:", jsonErr)
    }
    accessToken,ok := res["access_token"].(string)
    if !ok {
        fmt.Println("access token not found in response")
    }
    return accessToken
}


func main() {
        client := http.Client{}

        getPlaylists := getRequest{
            targetEndPoint: "https://api.spotify.com/v1/users/",
            authorisation: fmt.Sprintf("Bearer %s",getToken(clientId,clientSecret,client)),
            variable: "f5adff",
            value: "playlists",
        }


    resp2 := sendGetRequest(&getPlaylists,client)
    fmt.Println(resp2)
}
