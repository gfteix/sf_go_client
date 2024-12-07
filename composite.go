package main

import (
	"encoding/json"
	"errors"
)

type SubRequest struct {
	Method      string
	URL         string
	ReferenceID string
	Body        map[string]interface{}
	HTTPHeaders map[string]string
}

type CompositeError struct {
	StatusCode string   `json:"statusCode"`
	Message    string   `json:"message"`
	Fields     []string `json:"fields"`
}

/*
type ErrorBody []struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

type SuccessBody struct {
	Id      string   `json:"id"`
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}
*/

type SubResponse struct {
	Body           interface{}       `json:"body"`
	HTTPHeaders    map[string]string `json:"httpHeaders"`
	HTTPStatusCode int               `json:"httpStatusCode"`
	ReferenceID    string            `json:"referenceId"`
}

type CompositeResponse struct {
	CompositeResponse []SubResponse `json:"compositeResponse"`
}

type CompositeInput struct {
	client    SalesforceClient
	allOrNone bool
	requests  []SubRequest
}

type CollectionsRecord struct {
	sObjectType string
	fields      map[string]interface{}
}

type CollectonsProps struct {
	client    SalesforceClient
	allOrNone bool
	records   []CollectionsRecord
}

type CollectionsResponse []struct {
	Id      string          `json:"id"`
	Success bool            `json:"success"`
	Errors  []ErrorResponse `json:"errors"`
}

func Composite(i CompositeInput) (CompositeResponse, error) {
	body := make(map[string]interface{})

	body["allOrNone"] = i.allOrNone
	body["compositeRequests"] = i.requests

	result, err := i.client.fetch(FetchRequest{
		method: "POST",
		path:   "/composite",
		body:   body,
	})

	if err != nil {
		return CompositeResponse{}, err
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			return CompositeResponse{}, err
		}

		return CompositeResponse{}, errors.New(response[0].Message)
	}

	var response CompositeResponse

	if err := json.Unmarshal(result.body, &response); err != nil {
		return CompositeResponse{}, err
	}

	return response, nil
}

/*
Use a POST request with sObject Collections to add up to 200 records, returning a list of SaveResult objects. You can choose whether to roll back the entire request when an error occurs.

The list can contain up to 200 objects.
The list can contain objects of different types, including custom objects.
Each object must contain an attributes map. The map must contain a value for type.
*/
func Collections(props CollectonsProps) (CollectionsResponse, error) {
	body := make(map[string]interface{})

	body["allOrNone"] = props.allOrNone
	body["records"] = props.records

	/*
		var records []interface

		for _, v := range props.records {
			{
				// map to proper payload, define the attributes url and the fields
			}
		}*/

	result, err := props.client.fetch(FetchRequest{
		method: "POST",
		path:   "/composite/sobjects",
		body:   body,
	})

	if err != nil {
		return CollectionsResponse{}, err
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			return CollectionsResponse{}, err
		}

		return CollectionsResponse{}, errors.New(response[0].Message)
	}

	var response CollectionsResponse

	if err := json.Unmarshal(result.body, &response); err != nil {
		return CollectionsResponse{}, err
	}

	return response, nil
}
