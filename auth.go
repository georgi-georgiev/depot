package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Auth struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

var authToken *Auth

func GetAuthotization(address string, realm string) *Auth {
	if authToken == nil {
		endpoint := address + "auth/realms/" + realm + "/protocol/openid-connect/token"
		data := url.Values{}
		data.Set("client_id", "a1")
		data.Add("grant_type", "client_credentials")
		data.Add("client_secret", "1546156d-bd33-41f4-ac38-5e688d071bcf")
		data.Add("scope", "openid")

		req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)

		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &authToken)
		if err != nil {
			panic(err)
		}
	}

	return authToken
}
