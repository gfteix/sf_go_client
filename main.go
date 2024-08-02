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

	newAccount := make(map[string]interface{})
	newAccount["Name"] = "New Account"

	id, err := Create(client, "Account", newAccount)

	if err != nil {
		log.Panicf("Error while creating new account %v", err)
	}

	log.Printf("Account created: %v", id)

	query := fmt.Sprintf("SELECT Id, Name FROM Account WHERE Id = '%v'", id)

	result, err := Query(client, query)

	if err != nil {
		log.Panic("Failed to query Salesforce")
	}

	for _, record := range result {
		name := record["Name"]
		log.Printf("%v", name)

		id := record["Id"]
		log.Printf("%v", id)
	}
}
