package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
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

func Create(client *SalesforceClient, objectType string, body map[string]interface{}) (string, error) {
	result, statusCode := client.fetch(FetchProps{
		body:   body,
		method: "POST",
		url:    fmt.Sprintf("%v/services/data/v%v/sobjects/%v", c.orgUrl, c.apiVersion, objectType),
	})

	fmt.Printf("status code: %v", statusCode)

	var response CreateResponse

	if err := json.Unmarshal(result, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if !response.Success {
		return "", errors.New(strings.Join(response.Errors, "-"))
	}
	return response.Id, nil

}

func Update(client *SalesforceClient, objectType string, recordId string, body map[string]interface{}) error {
	result, statusCode := client.fetch(FetchProps{
		body:   body,
		method: "PATch",
		url:    fmt.Sprintf("%v/services/data/v%v/sobjects/%v/%v", c.orgUrl, c.apiVersion, objectType, recordId),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		return errors.New(response.Message)
	}

	return nil
}

func Delete(client *SalesforceClient, objectType string, recordId string) error {
	result, statusCode := client.fetch(FetchProps{
		method: "DELETE",
		url:    fmt.Sprintf("%v/services/data/v%v/sobjects/%v/%v", c.orgUrl, c.apiVersion, objectType, recordId),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		return errors.New(response.Message)
	}

	var response ErrorResponse

	if err := json.Unmarshal(result, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return errors.New(response.Message)
}

func Query(client *SalesforceClient, query string) ([]map[string]interface{}, error) {
	result, statusCode := client.fetch(FetchProps{
		method: "GET",
		url:    fmt.Sprintf("%v/services/data/v%v/q=%v", c.orgUrl, c.apiVersion, url.QueryEscape(query)),
	})

	fmt.Printf("status code: %v", statusCode)

	if *statusCode > 299 {
		var response ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		return nil, errors.New(response.Message)
	}

	var response []map[string]interface{}

	if err := json.Unmarshal(result, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return response, nil
}
