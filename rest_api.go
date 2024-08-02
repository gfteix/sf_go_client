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
	Done      bool                     `json:"done"`
	TotalSize int                      `json:"totalSize"`
	Records   []map[string]interface{} `json:"records"`
}

func Create(client *SalesforceClient, objectType string, body map[string]interface{}) (string, error) {
	result, statusCode := client.fetch(FetchProps{
		body:   body,
		method: "POST",
		path:   fmt.Sprintf("/sobjects/%v/", objectType),
	})

	if statusCode == nil {
		return "", errors.New("unable to fetch data")
	}

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
		path:   fmt.Sprintf("/sobjects/%v/%v", objectType, recordId),
	})

	if statusCode == nil {
		return errors.New("unable to fetch data")
	}

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
		path:   fmt.Sprintf("/sobjects/%v/%v", objectType, recordId),
	})

	if statusCode == nil {
		return errors.New("unable to fetch data")
	}
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
		path:   "/query/?q=" + url.QueryEscape(query),
	})

	if statusCode == nil {
		return nil, errors.New("unable to fetch data")
	}

	if *statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			return nil, err
		}

		return nil, errors.New(response[0].Message)
	}

	var response QueryResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response.Records, nil
}
