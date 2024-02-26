#!/bin/bash

targetUrl=$1
echo $targetUrl
accessToken=$(curl -X POST "https://accounts.spotify.com/api/token"\
                -H "Content-Type: application/x-www-form-urlencoded"\
                -d "grant_type=client_credentials&client_id=1b0ac2b304e941d9890dc016171c2226&client_secret=dd8f644ef4074f7f82daca80487818b6" | jq ".access_token"
            )


privToken="BQBTu-lL1jTXka4pftzc7PbrDVQm8L1JETSMdx7oP9ren9MSpirE7e6b34vOoj-rH7s3e7sDN1UVPHwwXlrPu8pLjMjm95ZgLRfFOfghu70ot2MXklYzxjZfUkZH1UG6TY54a53LRkvrP-mUiUqOnwB7xMRTyXGjVTmVqnTdA45lEca5NjEt9AtCWWW0nQK618Vhe0z_uX0z7ldbpyIemf4P1O7Vx1TrLooAAKivn7ktyoGEB0icJYsyCm5ob6KROU_sgxBncCtbCHfJIYNf29KSivssv8eiXdevLl27MhXlomjF"
echo $accessToken
curl --request GET --url $targetUrl --header 'Authorization: Bearer $privToken'
