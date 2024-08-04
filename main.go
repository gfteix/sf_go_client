package main

import (
	"fmt"
	"log"
	"reflect"

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

	query := fmt.Sprintf("SELECT Id, Name, Parent.Name, CreatedDate FROM Account WHERE Id = '%v'", id)
	result, err := Query(client, query)

	if err != nil {
		log.Panic("Failed to query Salesforce")
	}

	printResult(result)
}

func printResult(result []map[string]interface{}) {
	for index, record := range result {
		log.Printf("Index %v", index)

		for key, value := range record {
			log.Printf("%s: %v", key, value)
			log.Printf("Type: %v", reflect.TypeOf(value))
		}
	}
}
