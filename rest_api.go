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
	result, statusCode := props.client.fetch(FetchProps{
		body:   props.body,
		method: "POST",
		path:   fmt.Sprintf("/sobjects/%v/", props.objectType),
	})

	if statusCode == nil {
		return "", errors.New("unable to fetch data")
	}

	if *statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return "", err
		}

		return "", errors.New(response[0].Message)
	}

	var response CreateResponse

	if err := json.Unmarshal(result, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return "", err
	}

	return response.Id, nil

}

func Update(props UpdateProps) error {
	result, statusCode := props.client.fetch(FetchProps{
		body:   props.body,
		method: "PATCH",
		path:   fmt.Sprintf("/sobjects/%v/%v", props.objectType, props.recordId),
	})

	if statusCode == nil {
		return errors.New("unable to fetch data")
	}

	if *statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return err
		}

		return errors.New(response[0].Message)
	}

	return nil
}

func Delete(props DeleteProps) error {
	result, statusCode := props.client.fetch(FetchProps{
		method: "DELETE",
		path:   fmt.Sprintf("/sobjects/%v/%v", props.objectType, props.recordId),
	})

	if statusCode == nil {
		return errors.New("unable to fetch data")
	}

	if *statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
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

	if props.fields != nil && len(props.fields) > 0 {
		fields := strings.Join(props.fields, ",")
		path = path + "?fields=" + fields
	}

	result, statusCode := props.client.fetch(FetchProps{
		method: "GET",
		path:   path,
	})

	if statusCode == nil {
		return nil, errors.New("unable to fetch data")
	}

	if *statusCode > 299 {
		var response []ErrorResponse

		if err := json.Unmarshal(result, &response); err != nil {
			log.Println("Failed to unmarshal result")
			return nil, err
		}

		return nil, errors.New(response[0].Message)
	}

	var response Record

	if err := json.Unmarshal(result, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return nil, err
	}

	return response, nil
}

func Query(client *SalesforceClient, query string) ([]Record, error) {
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
			log.Println("Failed to unmarshal result")
			return nil, err
		}

		return nil, errors.New(response[0].Message)
	}

	var response QueryResponse

	if err := json.Unmarshal(result, &response); err != nil {
		log.Println("Failed to unmarshal result")
		return nil, err
	}

	return response.Records, nil
}
