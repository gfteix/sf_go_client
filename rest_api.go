package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type CreateResponse struct {
	Id      string   `json:"id"`
	Errors  []string `json:"errors"`
	Success bool     `json:"success"`
}

type QueryResponse struct {
	Done      bool     `json:"done"`
	TotalSize int      `json:"totalSize"`
	Records   []Record `json:"records"`
}

type CreateProps struct {
	client     *SalesforceClient
	objectType string
	body       map[string]interface{}
}

type UpdateProps struct {
	client     *SalesforceClient
	objectType string
	recordId   string
	body       map[string]interface{}
}

type DeleteProps struct {
	client     *SalesforceClient
	objectType string
	recordId   string
}

type GetByIdProps struct {
	client     *SalesforceClient
	objectType string
	recordId   string
	fields     []string // optional
}

func Create(props CreateProps) (string, error) {
	result, err := props.client.fetch(FetchRequest{
		body:   props.body,
		method: "POST",
		path:   fmt.Sprintf("/sobjects/%v/", props.objectType),
	})

	if err != nil {
		return "", err
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return "", err
		}

		return "", errors.New(response[0].Message)
	}

	var response CreateResponse

	if err := json.Unmarshal(result.body, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return "", err
	}

	return response.Id, nil

}

func Update(props UpdateProps) error {
	response := FetchResponse{}

	response, error := props.client.fetch(FetchRequest{
		body:   props.body,
		method: "PATCH",
		path:   fmt.Sprintf("/sobjects/%v/%v", props.objectType, props.recordId),
	})

	if error == nil {
		return error
	}

	if response.statusCode > 299 {
		var errorResponse []ErrorResponse

		if err := json.Unmarshal(response.body, &errorResponse); err != nil {
			log.Println("Failed to unmarshal result")
			return err
		}

		return errors.New(errorResponse[0].Message)
	}

	return nil
}

func Delete(props DeleteProps) error {
	result, err := props.client.fetch(FetchRequest{
		method: "DELETE",
		path:   fmt.Sprintf("/sobjects/%v/%v", props.objectType, props.recordId),
	})

	if err != nil {
		return errors.New("unable to fetch data")
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			return err
		}

		return errors.New(response[0].Message)
	}

	return nil
}

func GetByExternalId() (map[string]interface{}, error) {
	log.Panic("Not implemented")
	return nil, nil
}

func GetById(props GetByIdProps) (map[string]interface{}, error) {
	path := fmt.Sprintf("/sobjects/%v/%v/", props.objectType, props.recordId)

	if len(props.fields) > 0 {
		fields := strings.Join(props.fields, ",")
		path = path + "?fields=" + fields
	}

	result, error := props.client.fetch(FetchRequest{
		method: "GET",
		path:   path,
	})

	if error != nil {
		return nil, error
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return nil, err
		}

		return nil, errors.New(response[0].Message)
	}

	var response Record

	if err := json.Unmarshal(result.body, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return nil, err
	}

	return response, nil
}

func Query(client *SalesforceClient, query string) ([]Record, error) {
	result, error := client.fetch(FetchRequest{
		method: "GET",
		path:   "/query/?q=" + url.QueryEscape(query),
	})

	if error != nil {
		return nil, error
	}

	if result.statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result.body, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return nil, err
		}

		return nil, errors.New(response[0].Message)
	}

	var response QueryResponse

	if err := json.Unmarshal(result.body, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return nil, err
	}

	return response.Records, nil
}
