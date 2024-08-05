# SF GO Client

A simple Salesforce API client.

## Initiating the client

For now, only the password auth is supported.

Make sure the following environment variables are defined

```
USERNAME=
PASSWORD=
CLIENT_ID=
CLIENT_SECRET=
ORG_URL=
API_VERSION=
```


Call the NewSalesforceClient to retrieve a new client

```
    client := NewSalesforceClient()
```

## Rest API

### Creating a record

```GO
	newAccount := make(map[string]interface{})
	newAccount["Name"] = "New Account"

	id, err := Create(CreateProps{
		client:     client,
		objectType: "Account",
		body:       newAccount,
	})

	if err != nil {
		log.Printf("Error while creating new account %v", err)
	}
```

### Updating a record

```GO
	account := make(map[string]interface{})
	account["Name"] = "Updated Value"

	err := Update(UpdateProps{
		client:     client,
		objectType: "Account",
		recordId: 	"recordId",
		body:       account,
	})

	if err != nil {
		log.Printf("Error while updating account %v", err)
	}
```

### Deleting a record

```GO
	err := Delete(DeleteProps{
		client: client, 
		objectType: "Account", 
		recordId: "Salesforce RecordId",
	})

	if err != nil {
		log.Printf("Error while deleting account %v", err)
	}
```

### Query

```GO
	query := fmt.Sprintf("SELECT Id, Name, Parent.Name, CreatedDate FROM Account WHERE Id = '%v'", id)
	result, err := Query(client, query)

	if err != nil {
		log.Panic("Failed to query Salesforce")
	}

	// Priting the result
	
	for index, record := range result {
		log.Printf("Index %v", index)
	
		for key, value := range record {
			log.Printf("%s: %v", key, value)
			log.Printf("Type: %v", reflect.TypeOf(value))
		}
	}

```

## Metadata API (TODO)


## Composite API (TODO)
