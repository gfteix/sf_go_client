package main

import (
	"fmt"
	"log"
)

func main() {
	client := NewSalesforceClient()

	newAccount := make(map[string]interface{})
	newAccount["Name"] = "New Account"

	id, err := Create(client, "Account", newAccount)

	if err != nil {
		log.Printf("Error while creating new account %v", err)
	} else {
		fmt.Printf("Account created: %v", id)
	}

	result, err := Query(client, "SELECT Id FROM Account LIMIT 1")

	if err != nil {
		log.Printf("Failed to query Salesforce")
	} else {
		log.Printf("Query Result: %v", result)
	}
}
