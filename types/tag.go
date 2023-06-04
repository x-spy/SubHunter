package types

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func getTagUrl(token string) string {
	domain := "https://sub.tagsub.net/api/v1/client/subscribe?token="
	return domain + token
}

// Send request and return state
func request(client *http.Client) (*http.Response, string, error) {

	token := tokenGen()
	tagUrl := getTagUrl(token)
	fmt.Println("Token: " + token)

	resp, err := client.Get(tagUrl)
	if err != nil {
		fmt.Println("Error getting subscribe information, please check your internet connection.", err)
		return nil, "", err
	}
	return resp, token, nil

}

// Token generator
func tokenGen() string {
	const letterString = "0123456789abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, 32)
	for i := range token {
		token[i] = letterString[rand.Intn(len(letterString))]
	}
	return string(token)
}

func TagStart(threads int) {
	client := &http.Client{Timeout: 10 * time.Second}

	var wg sync.WaitGroup
	sem := make(chan bool, threads)

	for {
		wg.Add(1)
		sem <- true
		go func(c *http.Client) {

			defer wg.Done()
			defer func() { <-sem }()

			resp, token, err := request(client)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println("Failed to close response body, interrupting...")
					return
				}
			}(resp.Body)

			if err != nil {
				//return
			}
			if resp.StatusCode == http.StatusOK {
				fmt.Println("Available token: " + token + "\n" + "URL: " + getTagUrl(token))
				return
			}
		}(client)
	}
}
