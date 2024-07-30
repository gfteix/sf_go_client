package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	client := NewSalesforceClient()

	/*
		newAccount := make(map[string]interface{})
		newAccount["Name"] = "New Account"

		id, err := Create(client, "Account", newAccount)

		if err != nil {
			log.Printf("Error while creating new account %v", err)
		} else {
			log.Printf("Account created: %v", id)
		}
	*/

	query := fmt.Sprintf("SELECT Id, Name FROM Account WHERE Id = '%v'", "001aj00000RjWfKAAV")

	result, err := Query(client, query)

	if err != nil {
		log.Printf("Failed to query Salesforce")
	} else {
		log.Printf("Query Result: %v", result)
	}
}
