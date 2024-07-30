package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type CreateResponse struct {
	Id      string   `json:"id"`
	Errors  []string `json:"errors"`
	Success bool     `json:"success"`
}

type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
	Fields    string `json:"fields"`
}

type QueryResponse struct {
	Done      string                   `json:"done"`
	TotalSize int                      `json:"totalSize"`
	Records   []map[string]interface{} `json:"records"`
}

func Create(client *SalesforceClient, objectType string, body map[string]interface{}) (string, error) {
	result, statusCode := client.fetch(FetchProps{
		body:   body,
		method: "POST",
		url:    fmt.Sprintf("%v/sobjects/%v/", client.apiUrl, objectType),
	})

	fmt.Printf("status code: %v", *statusCode)

	if *statusCode > 299 {
		var response []map[string]interface{}

		if err := json.Unmarshal(result, &response); err != nil {
			return "", err
		}

		message := response[0]["message"].(string)

		return "", errors.New(message)
	}

	var response CreateResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return "", err
	}

	return response.Id, nil

}

func Update(client *SalesforceClient, objectType string, recordId string, body map[string]interface{}) error {
	result, statusCode := client.fetch(FetchProps{
		body:   body,
		method: "PATch",
		url:    fmt.Sprintf("%v/sobjects/%v/%v", client.apiUrl, objectType, recordId),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			return err
		}

		return errors.New(response.Message)
	}

	return nil
}

func Delete(client *SalesforceClient, objectType string, recordId string) error {
	result, statusCode := client.fetch(FetchProps{
		method: "DELETE",
		url:    fmt.Sprintf("%v/sobjects/%v/%v", client.apiUrl, objectType, recordId),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			return err
		}

		return errors.New(response.Message)
	}

	var response ErrorResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return err
	}

	return errors.New(response.Message)
}

func Query(client *SalesforceClient, query string) ([]map[string]interface{}, error) {
	result, statusCode := client.fetch(FetchProps{
		method: "GET",
		url:    fmt.Sprintf("%v/query/?q=%v", client.apiUrl, url.QueryEscape(query)),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			return nil, err
		}

		return nil, errors.New(response.Message)
	}

	var response QueryResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response.Records, nil
}
