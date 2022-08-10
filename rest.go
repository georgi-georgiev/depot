package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SetTokenToHeader(req *http.Request, address string, realm string) {
	auth := GetAuthotization(address, realm)

	req.Header.Add("Authorization", "Bearer "+auth.AccessToken)

	log.Printf("authorization token: %s", auth.AccessToken)
}

func PrepareRequest(req *http.Request, useToken bool, address string, realm string) {
	if useToken {
		SetTokenToHeader(req, address, realm)
	}
	req.Header.Add("Content-Type", "application/json")
}

func PrepareBody(file string) *bytes.Buffer {
	json := LoadJson(file)
	return bytes.NewBuffer(json)
}

func ValidateResponse(resp *http.Response) []byte {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("response body: %s", string(bodyBytes))

	if resp.StatusCode == http.StatusOK {
		return bodyBytes
	} else {
		message := fmt.Errorf("expected to get status %d but instead got %d", http.StatusOK, resp.StatusCode)
		panic(message)
	}
}

func Execute(req *http.Request, useToken bool, address string, realm string) []byte {
	PrepareRequest(req, useToken, address, realm)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return ValidateResponse(resp)
}

func Post(address string, file string, authAddress string, realm string) []byte {
	return PostWithToken(address, file, false, authAddress, realm)
}

func PostWithToken(address string, file string, useToken bool, authAddress string, realm string) []byte {
	log.Printf("request to address: %s", address)

	body := PrepareBody(file)

	req, err := http.NewRequest(http.MethodPost, address, body)
	if err != nil {
		panic(err)
	}

	return Execute(req, useToken, authAddress, realm)
}

func Get(address string, query map[string]string, authAddress string, realm string) []byte {
	return GetWithToken(address, query, false, authAddress, realm)
}

func GetWithToken(address string, query map[string]string, useToken bool, authAddress string, realm string) []byte {
	log.Printf("request to address: %s\n", address)

	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		panic(err)
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}

		rawQuery := q.Encode()
		req.URL.RawQuery = rawQuery

		log.Printf("raw query: %s", rawQuery)
		log.Printf("raw adress: %s", req.URL)
	}

	return Execute(req, useToken, authAddress, realm)
}
