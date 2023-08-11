package transmissionrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const csrfHeader = "X-Transmission-Session-Id"

type requestPayload struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
	Tag       int         `json:"tag,omitempty"`
}

type answerPayload struct {
	Arguments interface{} `json:"arguments"`
	Result    string      `json:"result"`
	Tag       *int        `json:"tag"`
}

func (c *Client) rpcCall(ctx context.Context, method string, arguments interface{}, result interface{}) (err error) {
	return c.request(ctx, method, arguments, result, true)
}

func (c *Client) request(ctx context.Context, method string, arguments interface{}, result interface{}, retry bool) (err error) {
	// Let's avoid crashing if not instanciated properly
	if c.http == nil {
		err = errors.New("this controller is not initialized, please use the New() function")
		return
	}
	// Prepare request payload
	rq := requestPayload{
		Method:    method,
		Arguments: arguments,
		Tag:       c.getRandomTag(),
	}
	rqJSON, err := json.Marshal(rq)
	if err != nil {
		err = fmt.Errorf("failed to marshal request payload: %w", err)
		return
	}
	// Build the request
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "POST", c.endpoint.String(), bytes.NewBuffer(rqJSON)); err != nil {
		err = fmt.Errorf("can't prepare request for '%s' method: %w", method, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set(csrfHeader, c.getSessionID())
	// Execute request
	var resp *http.Response
	if resp, err = c.http.Do(req); err != nil {
		err = fmt.Errorf("failed to execute HTTP request: %w", err)
		return
	}
	defer resp.Body.Close()
	// Is the CRSF token invalid ?
	if resp.StatusCode == http.StatusConflict {
		// Recover new token and save it
		c.updateSessionID(resp.Header.Get(csrfHeader))
		// Retry request if first try
		if retry {
			return c.request(ctx, method, arguments, result, false)
		}
		err = errors.New("CSRF token invalid 2 times in a row: stopping to avoid infinite loop")
		return
	}
	// Is request successful ?
	if resp.StatusCode != 200 {
		err = HTTPStatusCode(resp.StatusCode)
		return
	}
	// Decode body
	answer := answerPayload{
		Arguments: result,
	}
	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		err = fmt.Errorf("can't unmarshal request answer body: %w", err)
		return
	}
	// Final checks
	if answer.Tag == nil {
		err = errors.New("http answer does not have a tag within it's payload")
		return
	}
	if *answer.Tag != rq.Tag {
		err = errors.New("http request tag and answer payload tag do not match")
		return
	}
	if answer.Result != "success" {
		err = fmt.Errorf("http request ok but payload does not indicate success: %s", answer.Result)
		return
	}
	// All good
	return
}

// HTTPStatusCode is a custom error type for HTTP errors
type HTTPStatusCode int

func (hsc HTTPStatusCode) Error() string {
	text := http.StatusText(int(hsc))
	if text != "" {
		text = ": " + text
	}
	return fmt.Sprintf("HTTP error %d%s", hsc, text)
}
