package fcm

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseFcmResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"k2d2e3r4",
					"registration_id": "t5y6e9d8o9",
					"error": ""
				}]
			}`)
	}))

	defer server.Close()

	data := `{"data":{"msg":"hello","sum":"New year"},"priority":"high","registration_ids":["token 1","token 2","token 3"]}`
	body := bytes.NewBuffer([]byte(data))

	res, err := http.Post(server.URL, "application/json", body)
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
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, `{
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
	}))

	defer server.Close()

	// Set token
	token := "token1"

	res, err := http.Get(server.URL + fmt.Sprintf("?token=%s", token))
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
