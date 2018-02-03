package TransmissionRPC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const csrfHeader = "X-Transmission-Session-Id"

type requestPayload struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
	Tag       int         `json:"tag,omitempty"`
}

func (c *Controller) request(method string, arguments interface{}) (err error) {
	// Prepare payload
	var buff bytes.Buffer
	err = json.NewEncoder(&buff).Encode(&requestPayload{
		Method:    method,
		Arguments: arguments,
		// Tag:       tag,
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
	req.Header.Set(csrfHeader, c.sessionID)
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
		c.sessionID = resp.Header.Get(csrfHeader)
		fmt.Printf("DEBUG updated CSRF token to: %s\n", c.sessionID)
		// Retry request
		if resp, err = c.httpC.Do(req); err != nil {
			err = fmt.Errorf("request error: %v", err)
			return
		}
	}
	// Is request successfull ?
	if resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP error %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}
	// Decode body
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	return
	// var answer interface{}
	// if err = json.NewDecoder(resp.Body).Decode(answer); err != nil {
	// 	err = fmt.Errorf("can't unmarshall request answer body: %v", err)
	// 	return
	// }
	// // tmp
	// fmt.Println("wahou")
	// fmt.Println(answer)
	// return
}
