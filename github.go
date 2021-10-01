package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	githubURL = "https://api.github.com/app/installations"
)

//  curl -i -X POST -H "Authorization: Bearer $TOKEN" -H "Accept: " https://api.github.com/app/installations/19557896/access_tokens

// {
//   "token": "asdf",
//   "expires_at": "2021-09-17T14:00:44Z",
//   "permissions": {
//     "contents": "read",
//     "metadata": "read",
//     "pull_requests": "write"
//   },
//   "repository_selection": "selected"
// }

func GetInstallationToken(token string) (*string, error) {

	u := strings.Join([]string{githubURL, appInstId, "access_tokens"}, "/")

	var resBody struct {
		Token       string    `json:"token"`
		ExpiresAt   time.Time `json:"expires_at"`
		Permissions struct {
			Contents     string `json:"contents"`
			Metadata     string `json:"metadata"`
			PullRequests string `json:"pull_requests"`
		} `json:"permissions"`
		RepositorySelection string `json:"repository_selection"`
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		//TODO handle error
		fmt.Println("something went wrong", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request Error:%s\n", err)
		return nil, err
	}

	b, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		fmt.Println(res.StatusCode)
		log.Println(string(b))
		return nil, errors.New("Invalid response code")
	}

	if err := json.Unmarshal(b, &resBody); err != nil {
		log.Printf("Problem unmarshalling err:%s\n", err)
	}

	return &resBody.Token, nil
}
