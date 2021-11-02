package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Message string `json:"error_message"`
	// Code    int    `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf(e.Message)
}

type accountMessage struct {
	// Code int `json:"code"`
	Data interface{} `json:"data"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

const ErrUnknowType = "unknown error"

func NewClient(basePath string, hc *http.Client) *Client {
	if hc == nil {
		hc = &http.Client{
			Timeout: time.Minute,
		}
	}
	if len(basePath) == 0 {
		basePath = "http://localhost:8080/v1"
	}
	return &Client{
		baseURL:    basePath,
		httpClient: hc,
	}
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusCreated {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		// return errors.New(ErrUnknowType)
		errRes = ErrorResponse{Message: ErrUnknowType}
		return &errRes

		// return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := accountMessage{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
