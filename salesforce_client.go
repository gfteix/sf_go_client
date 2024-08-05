package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type SalesforceClient struct {
	apiVersion   int
	clientId     string
	clientSecret string
	password     string
	username     string
	orgUrl       string
	apiUrl       string
	token        *string
}

type FetchProps struct {
	body   map[string]interface{}
	method string
	path   string
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewSalesforceClient() *SalesforceClient {
	versionFromEnv := getEnv("API_VERSION")
	var version int

	if versionFromEnv != "" {
		v, err := strconv.ParseInt(versionFromEnv, 10, 0)

		if err != nil {
			log.Panicf("invalid version %v, %v", v, err)
		}

		version = int(v)
	} else {
		version = 61
	}

	c := &SalesforceClient{
		apiVersion:   version,
		clientId:     getEnv("CLIENT_ID"),
		clientSecret: getEnv("CLIENT_SECRET"),
		password:     getEnv("PASSWORD"),
		username:     getEnv("USERNAME"),
		orgUrl:       getEnv("ORG_URL"),
	}

	c.apiUrl = fmt.Sprintf("%v/services/data/v%v.0", c.orgUrl, c.apiVersion)

	return c
}

func (c *SalesforceClient) fetch(props FetchProps) ([]byte, *int) {
	var bufferBody io.Reader

	if props.body != nil {
		encodedBody, err := json.Marshal(props.body)

		if err != nil {
			log.Printf("Error converting body to JSON: %v", err)
			return nil, nil
		}

		bufferBody = bytes.NewBuffer(encodedBody)
	}

	client := &http.Client{}

	reqUrl := c.apiUrl + props.path
	req, err := http.NewRequest(props.method, reqUrl, bufferBody)

	if err != nil {
		log.Printf("error on http.NewRequest: %v", err)
		return nil, nil
	}

	token, err := c.getToken()

	if err != nil {
		log.Printf("error getting token %v", err)
		return nil, nil
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("error on client.Do: %v", err)
		return nil, nil
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("error on io.ReadAll: %v", err)
		return nil, &resp.StatusCode
	}

	fmt.Printf("status code: %v\n", resp.StatusCode)

	return respBody, &resp.StatusCode
}

func (c *SalesforceClient) getToken() (string, error) {
	if c.token != nil {
		return *c.token, nil
	}

	requestUrl := fmt.Sprintf("%s/services/oauth2/token", c.orgUrl)
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", c.clientId)
	data.Set("client_secret", c.clientSecret)
	data.Set("username", c.username)
	data.Set("password", c.password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(data.Encode()))

	if err != nil {
		log.Printf("error on http.NewRequest: %v", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("error on client.Do: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var result TokenError

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}
		log.Printf("Error getting token: %v", result.ErrorDescription)

		return "", errors.New("Error getting the token: " + result.ErrorDescription)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response body: %v", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("access_token not found in response")
	}

	c.token = &token

	return token, nil
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ""
}
