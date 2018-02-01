# go-fcm : FCM library in Go

[![CircleCI](https://circleci.com/gh/douglasmakey/go-fcm.svg?style=svg)](https://circleci.com/gh/douglasmakey/go-fcm)
[![Codecov]](https://codecov.io/gh/douglasmakey/go-fcm)

Firebase Cloud Messaging ( FCM ) Library using golang ( Go )

This library uses HTTP/JSON Firebase Cloud Messaging connection server protocol

## Usage

```
go get github.com/douglasmakey/go-fcm
```

## Docs

####  Firebase Cloud Messaging HTTP Protocol Specs
```
https://firebase.google.com/docs/cloud-messaging/http-server-ref
```

#### Firebase Cloud Messaging Developer docs
```
https://firebase.google.com/docs/cloud-messaging/
```

## Example

```go
package main

import (
	"log"
	"github.com/douglasmakey/go-fcm"
	)

func main() {
	// init client
	client := fcm.NewClient("ApiKey")
	
	// You can use your HTTPClient 
	//client.SetHTTPClient(client)
	
	data := map[string]interface{}{
		"message": "From Go-FCM",
		"details": map[string]string{
			"name": "Name",
			"user": "Admin",
			"thing": "none",
		},
	}
	
	// You can use PushMultiple or PushSingle
	client.PushMultiple([]string{"token 1", "token 2"}, data)
	//client.PushSingle("token 1", data)
	
	// registrationIds remove and return map of invalid tokens
	badRegistrations := client.CleanRegistrationIds()
	log.Println(badRegistrations) 
	
	status, err := client.Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	
	log.Println(status.Results)
}

```
[Codecov]: https://codecov.io/gh/douglasmakey/go-fcm/branch/master/graph/badge.svg