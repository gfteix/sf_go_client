# SF GO Client

A simple client to interact with Salesforce API.



## Initiating the client

For now, only the password auth is supported.

Make sure the following variables are defined:

```
USERNAME=
PASSWORD=
CLIENT_ID=
CLIENT_SECRET=
ORG_URL=
API_VERSION= // optional
```


Call the NewSalesforceClient to retrieve a new client

```
    client := NewSalesforceClient()

```

Then it is just a matter of passing the client to the helper functions.

## Metadata API (TODO)

## Rest API

*Creating a record*

```	
	newAccount := make(map[string]interface{})
	newAccount["Name"] = "New Account"

	id, err := Create(client, "Account", newAccount)

	if err != nil {
		log.Printf("Error while creating new account %v", err)
	}
```

*Updating a record*

```
	account := make(map[string]interface{})
	account["Name"] = "Updated Value"

	err := Create(client, "Account", "Salesforce recordId", account)

	if err != nil {
		log.Printf("Error while updating account %v", err)
	}
```

*Deleting a record*

```
	err := Delete(client, "Account", ""Salesforce RecordId")

	if err != nil {
		log.Printf("Error while deleting account %v", err)
	}
```

*Query*

```

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
## Composite API (TODO)