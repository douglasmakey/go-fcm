package fcm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func parseFcmResponse(resp *http.Response) (*response, error) {
	// Defers
	defer resp.Body.Close()

	// Check statusCode from resp
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statusCode: %d error: %s", resp.StatusCode, resp.Status)
	}

	// Create response
	response := new(response)
	response.StatusCode = resp.StatusCode
	response.RetryAfter = resp.Header.Get("Retry-After")

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func parseTokenDetails(resp *http.Response) (*tokenDetails, error) {
	// Defers
	defer resp.Body.Close()

	// Create tokenDetails and decode
	tokenDetails := new(tokenDetails)
	tokenDetails.StatusCode = resp.StatusCode
	if err := json.NewDecoder(resp.Body).Decode(tokenDetails); err != nil {
		return nil, err
	}

	return tokenDetails, nil

}
