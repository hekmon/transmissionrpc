package TransmissionRPC

import (
	"bytes"
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
	Tag       int         `json:"tag"`
}

func (c *Controller) rpcCall(method string, arguments interface{}, result interface{}) (err error) {
	return c.request(method, arguments, result, true)
}

func (c *Controller) request(method string, arguments interface{}, result interface{}, retry bool) (err error) {
	// Prepare payload
	tag := c.rnd.Int()
	var buff bytes.Buffer
	err = json.NewEncoder(&buff).Encode(&requestPayload{
		Method:    method,
		Arguments: arguments,
		Tag:       tag,
	})
	if err != nil {
		err = fmt.Errorf("can't marshall JSON payload: %v", err)
		return
	}
	// Prepare the request
	var req *http.Request
	if req, err = http.NewRequest("POST", c.url, &buff); err != nil {
		err = fmt.Errorf("can't prepare request for '%s' method: %v", method, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set(csrfHeader, c.getSessionID())
	req.SetBasicAuth(c.user, c.password)
	// Execute request
	var resp *http.Response
	if resp, err = c.httpC.Do(req); err != nil {
		err = fmt.Errorf("request error: %v", err)
		return
	}
	defer resp.Body.Close()
	// Is the CRSF token invalid ?
	if resp.StatusCode == http.StatusConflict {
		// Recover new token
		c.updateSessionID(resp.Header.Get(csrfHeader))
		fmt.Printf("DEBUG updated CSRF token to: %s\n", c.getSessionID())
		// Retry request if first try
		if retry {
			return c.request(method, arguments, result, false)
		}
		err = errors.New("CSRF token invalid 2 times in a row: stopping to avoid infinite loop")
		return
	}
	// Is request successfull ?
	if resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP error %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}
	// Debug
	// {
	// 	var data []byte
	// 	data, err = ioutil.ReadAll(resp.Body)
	// 	fmt.Println(string(data))
	// 	return
	// }
	// Decode body
	answer := answerPayload{
		Arguments: result,
	}
	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		err = fmt.Errorf("can't unmarshall request answer body: %v", err)
		return
	}
	// Final checks
	if answer.Result != "success" {
		err = fmt.Errorf("http request ok but payload does not indicate success: %s", answer.Result)
		return
	}
	if answer.Tag != tag {
		err = errors.New("http request and answer payload tag do not match")
		return // not really needed but clean
	}
	// All good
	return
}

func (c *Controller) getSessionID() string {
	defer c.sessionIDAccess.RUnlock()
	c.sessionIDAccess.RLock()
	return c.sessionID
}

func (c *Controller) updateSessionID(newID string) {
	defer c.sessionIDAccess.Unlock()
	c.sessionIDAccess.Lock()
	c.sessionID = newID
}
