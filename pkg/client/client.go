package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/onspaceship/booster/pkg/config"

	"github.com/apex/log"
	"github.com/golang-jwt/jwt"
)

type Options = config.ClientOptions

type Client struct {
	*Options
}

func NewClient() *Client {
	options, err := config.NewClientOptions()
	if err != nil {
		log.WithError(err).Fatal("failed to configure API client")
	}

	return &Client{Options: options}
}

func (client *Client) Get(url string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	token, err := client.newToken()
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(data)

	return err
}

func (client *Client) Put(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	token, err := client.newToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (client *Client) newToken() (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	token.Claims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	}

	return token.SignedString(client.CoreJWTKey)
}
