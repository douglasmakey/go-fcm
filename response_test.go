package fcm

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"net/http"
)

func TestResponse_GetInvalidTokens(t *testing.T) {
	defer gock.Off()

	gock.New(apiFCM).
		Post("").
		Reply(http.StatusOK).
		JSON(`{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"1y2t4i224uy2b",
					"registration_id": "jfey12fugyuy12oijd",
					"error": ""
				}, {"error": "InvalidToken"}, {"error": "InvalidToken"}]
			}`)

	data := map[string]string{
		"body": "Test",
	}

	// Init client
	client := NewClient("key")
	client.PushMultiple([]string{"token 1"}, data)

	invalidTokens := []string{"token 2", "token 3"}
	client.AppendRegistrationIds(invalidTokens)

	status, err := client.Send()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if status.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", status.StatusCode)
	}

	if status.Success != 1 {
		t.Errorf("expected 1 got %d", status.Success)
	}


	if status.Failure != 0 {
		t.Errorf("expected 0 got %d", status.Failure)
	}

	badTokens := status.GetInvalidTokens()

	if len(badTokens) != 2 {
		t.Errorf("expected 2, got %d", len(badTokens))
	}

	for _, val := range badTokens {
		if val != "InvalidToken" {
			t.Errorf("expected InvalidToken, got %s", val)
		}
	}
}

