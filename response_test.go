package fcm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_GetInvalidTokens(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "key=test" {
			t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
		}
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"1y2t4i224uy2b",
					"registration_id": "jfey12fugyuy12oijd",
					"error": ""
				}, {"error": "InvalidToken"}, {"error": "InvalidToken"}]
			}`)
	}))

	defer server.Close()

	data := map[string]string{
		"body": "Test",
	}

	// Init client
	client := NewClient("test")
	client.ApiFCM = server.URL
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
