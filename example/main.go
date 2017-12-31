package main

import (
	"github.com/douglasmakey/go-fcm"
	"log"
)

func main() {
	// init client
	client := fcm.NewClient("ApiKey")

	// You can use your HTTPClient
	//client.SetHTTPClient(client)

	data := map[string]interface{}{
		"message": "From Go-FCM",
		"details": map[string]string{
			"name":  "Name",
			"user":  "Admin",
			"thing": "none",
		},
	}

	// You can use PushMultiple or PushSingle
	client.PushMultiple([]string{"token 1", "token 2"}, data)
	//client.PushSingle("token 1", data)

	// registrationIds remove and return a list of invalid tokens
	badRegistrations := client.CleanRegistrationIds()
	log.Println(badRegistrations)

	status, err := client.Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println(status.Results)

}
