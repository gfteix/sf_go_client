package main

type Record map[string]interface{}

type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
	Fields    string `json:"fields"`
}
