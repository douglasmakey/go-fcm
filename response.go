package fcm

type response struct {
	StatusCode          int
	Err                 string   `json:"error,omitempty"`
	Success             int      `json:"success"`
	MultiCastId         int64    `json:"multicast_id"`
	CanonicalIds        int      `json:"canonical_ids"`
	Failure             int      `json:"failure"`
	Results             []result `json:"results,omitempty"`
	MsgId               int64    `json:"message_id,omitempty"`
	RetryAfter          string   `json:"retry_after"`
	copyRegistrationIds []string
}

type result struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          string `json:"error"`
}

// GetInvalidTokens return list with tokens wrongs
func (r *response) GetInvalidTokens() map[string]string {
	tr := make(map[string]string)
	for index, val := range r.Results {
		if val.Error != "" {
			tr[r.copyRegistrationIds[index]] = val.Error
		}
	}

	return tr
}
