package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net"
	"time"
)

const (
	// Methods
	GET  = "GET"
	POST = "POST"

	// Define Urls
	defaultApiFCM = "https://fcm.googleapis.com/fcm/send"
	defaultApiIID = "https://iid.googleapis.com/iid/info/%s?details=true"

	// Set Max time to live for message
	maxTTL = 2419200

	// Priorities
	HighPriority   = "high"
	NormalPriority = "normal"
)

var (
	// Errors
	ErrDataIsEmpty  = errors.New("data is empty")
	ErrToManyRegIDs = errors.New("too many registrations ids")
)

type NotificationPayload struct {
	Title            string `json:"title,omitempty"`
	Body             string `json:"body,omitempty"`
	BodyLocKey       string `json:"body_loc_key,omitempty"`
	BodyLocArgs      string `json:"body_loc_args,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Color            string `json:"color,omitempty"`
	ClickAction      string `json:"click_action,omitempty"`
	TitleLocKey      string `json:"title_loc_key,omitempty"`
	TitleLocArgs     string `json:"title_loc_args,omitempty"`
	AndroidChannelID string `json:"android_channel_id,omitempty"`
}

type message struct {
	Data                  interface{}          `json:"data,omitempty"`
	To                    string               `json:"to,omitempty"`
	Notification          *NotificationPayload `json:"notification,omitempty"`
	Priority              string               `json:"priority,omitempty"`
	RegistrationIds       []string             `json:"registration_ids,omitempty"`
	MutableContent        bool                 `json:"mutable_content,omitempty"`
	Condition             string               `json:"condition,omitempty"`
	CollapseKey           string               `json:"collapse_key,omitempty"`
	ContentAvailable      bool                 `json:"content_available,omitempty"`
	RestrictedPackageName string               `json:"restricted_package_name,omitempty"`
	DryRun                bool                 `json:"dry_run,omitempty"`
	TimeToLive            int                  `json:"time_to_live,omitempty"`
}

type tokenDetails struct {
	Application      string `json:"application,omitempty"`
	Platform         string `json:"platform,omitempty"`
	AppSigner        string `json:"appSigner,omitempty"`
	AttestStatus     string `json:"attestStatus,omitempty"`
	AuthorizedEntity string `json:"authorizedEntity,omitempty"`
	ConnectionType   string `json:"connectionType,omitempty"`
	ConnectDate      string `json:"connectDate,omitempty"`
	StatusCode       int
	Error            string                                  `json:"error,omitempty"`
	Rel              map[string]map[string]map[string]string `json:"rel,omitempty"`
}

type Client struct {
	apiKey     string
	Message    *message
	clientHttp *http.Client
	ApiFCM     string
	ApiIID     string
}

// NewClient Create instance of client
func NewClient(key string) *Client {
	// Generate new client with apiKey
	client := new(Client)
	client.apiKey = key
	client.Message = &message{}

	// Create default HTTPClient
	c := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	client.clientHttp = c

	// Set default endpoints
	client.ApiFCM = defaultApiFCM
	client.ApiIID = defaultApiIID

	return client
}

// SetHTTPClient set specific HTTPClient
func (c *Client) SetHTTPClient(client *http.Client) {
	c.clientHttp = client
}

// SetData Set data for message
func (c *Client) SetData(d interface{}) {
	c.Message.Data = d
}

// SetMsgAndTo 'To' this parameter specifies the recipient of a message.
func (c *Client) PushSingle(to string, d interface{}) {
	c.SetData(d)
	c.Message.To = to
}

// SetMsgAndIds Set Message and ids for send
func (c *Client) PushMultiple(ids []string, d interface{}) {
	c.SetData(d)
	c.Message.RegistrationIds = ids
}

// AppendRegistrationIds Append your list of RegistrationIds
func (c *Client) AppendRegistrationIds(ids []string) {
	c.Message.RegistrationIds = append(c.Message.RegistrationIds, ids...)
}

// CleanRegistrationIds remove invalid token of RegistrationIds and return map of BadTokens
func (c *Client) CleanRegistrationIds() []string {
	var validTokens []string
	var badTokens []string
	for _, t := range c.Message.RegistrationIds {
		details, err := c.GetTokenDetails(t)
		if err == nil && details.Error == "" {
			validTokens = append(validTokens, t)
		} else {
			badTokens = append(badTokens, t)
		}
	}

	// Change RegistrationIds for validTokens
	c.Message.RegistrationIds = validTokens

	return badTokens
}

// GetTokenDetails get info about the token
func (c *Client) GetTokenDetails(t string) (*tokenDetails, error) {

	var url string
	if c.ApiIID == defaultApiIID {
		url = fmt.Sprintf(c.ApiIID, t)
	} else {
		url = c.ApiIID + fmt.Sprintf("?token=%s", t)
	}

	resp, err := c.doRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}

	details, err := parseTokenDetails(resp)
	if err != nil {
		return nil, err
	}

	return details, nil

}

// Send Validate and Send FCM message
func (c *Client) Send() (*response, error) {
	err := c.validateData()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(c.Message)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(POST, c.ApiFCM, b)
	if err != nil {
		return nil, err
	}

	response, err := parseFcmResponse(resp)
	if err != nil {
		return nil, err
	}

	// Copy registrations from Client to response
	response.copyRegistrationIds = c.Message.RegistrationIds

	return response, nil
}

// validateData return error if data is wrong
func (c *Client) validateData() error {
	// Data is empty
	if c.Message.Data == nil {
		return ErrDataIsEmpty
	}

	// Max token permit for FCM is 1000
	if len(c.Message.RegistrationIds) > 1000 {
		return ErrToManyRegIDs
	}

	// Validate Priority
	if c.Message.Priority != NormalPriority {
		c.Message.Priority = HighPriority
	}

	// Validate TimeToLive
	if c.Message.TimeToLive > maxTTL {
		c.Message.TimeToLive = maxTTL
	}

	return nil
}

// doRequest do request
func (c *Client) doRequest(m string, url string, data []byte) (*http.Response, error) {
	// Create request
	request, err := http.NewRequest(m, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Set headers
	request.Header.Set("Authorization", fmt.Sprintf("key=%v", c.apiKey))
	request.Header.Set("Content-Type", "application/json")

	// Execute requests
	resp, err := c.clientHttp.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
