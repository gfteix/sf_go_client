package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RecordAttributes struct {
	Type string `json:"type"`
}

type Record struct {
	Attributes RecordAttributes `json:"attributes"`
	// [key: string]: any -> how to map this to go type?
}

type CompositeBody struct {
	AllOrNone bool     `json:"allOrNone"`
	Records   []Record `json:"records"`
}

type CompositeError struct {
	StatusCode string   `json:"statusCode"`
	Message    string   `json:"message"`
	Fields     []string `json:"fields"`
}
type CompositeResponse struct { // array
	Success string           `json:"success"`
	Errors  []CompositeError `json:"CompositeError"`
}

type SalesforceClient struct {
	apiVersion   int
	clientId     string
	clientSecret string
	password     string
	username     string
	orgUrl       string
	token        *string
	tokenExpiry  time.Time
}

type FetchProps struct {
	body   map[string]interface{}
	method string
	url    string
}

func (c *SalesforceClient) fetch(props FetchProps) ([]byte, *int) {
	var bufferBody *bytes.Buffer = nil

	if props.body != nil {
		encodedBody, _ := json.Marshal(props.body)
		bufferBody = bytes.NewBuffer(encodedBody)
	}

	client := &http.Client{}
	req, err := http.NewRequest(props.method, props.url, bufferBody)

	if err != nil {
		log.Printf("error on http.NewRequest: %v", err)
		return nil, nil
	}

	token, err := c.getToken()

	if err != nil {
		log.Printf("error getting token %v", err)
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

	return respBody, &resp.StatusCode

}

func (c *SalesforceClient) getToken() (string, error) {
	if c.token != nil && time.Now().Before(c.tokenExpiry) {
		return *c.token, nil
	}

	url := fmt.Sprintf("%s/services/oauth2/token", c.orgUrl)
	body := make(map[string]string)

	body["grant_type"] = "password"
	body["client_id"] = c.clientId
	body["client_secret"] = c.clientSecret
	body["username"] = c.username
	body["password"] = c.password

	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, requestBody)

	if err != nil {
		log.Printf("error on http.NewRequest: %v", err)
		return "", err

	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json; charset=UTF-8")

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("error on client.Do: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("error on io.ReadAll: %v", err)
		return "", err
	}

	fmt.Println(string(respBody))

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response body: %v", err)
	}

	*c.token = result["access_token"].(string)
	c.tokenExpiry = result["expiry_time"].(time.Time) // ?

	return *c.token, nil
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ""
}

func NewSalesforceClient() *SalesforceClient {
	return &SalesforceClient{
		apiVersion:   61,
		clientId:     getEnv("CLIENT_ID"),
		clientSecret: getEnv("CLIENT_SECRET"),
		password:     getEnv("PASSWORD"),
		username:     getEnv("USERNAME"),
		orgUrl:       getEnv("ORG_URL"),
	}
}
