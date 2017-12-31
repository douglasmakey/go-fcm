package fcm

import (
	"bytes"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestParseFcmResponse(t *testing.T) {
	defer gock.Off()

	gock.New(apiFCM).
		Post("/").
		Reply(http.StatusOK).
		JSON(`{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)

	data := `{"data":{"msg":"hello","sum":"New year"},"priority":"high","registration_ids":["token 1","token 2","token 3"]}`
	body := bytes.NewBuffer([]byte(data))

	res, err := http.Post(apiFCM, "application/json", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response, err := parseFcmResponse(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", response.StatusCode)
	}
}

func TestParseTokenDetails(t *testing.T) {
	defer gock.Off()

	// Set token
	token := "token1"

	gock.New(apiIID).
		Get(token).
		Reply(http.StatusOK).
		JSON(`{
  				"application":"com.iid.example",
				"authorizedEntity":"123456782354",
			  	"platform":"Android",
			  	"attestStatus":"ROOTED",
			  	"appSigner":"1a2bc3d4e5",
			  	"connectionType":"WIFI",
			  	"connectDate":"2015-05-12",
			  	"rel":{
					"topics":{
				  		"topicname1":{"addDate":"2015-07-30"},
				  		"topicname2":{"addDate":"2015-07-30"},
				  		"topicname3":{"addDate":"2015-07-30"},
				  		"topicname4":{"addDate":"2015-07-30"}
						}
			  		}
	}`)

	res, err := http.Get(apiIID + token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	details, err := parseTokenDetails(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if details.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", details.StatusCode)
	}

}
