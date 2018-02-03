package TransmissionRPC

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const csrfHeader = "X-Transmission-Session-Id"

type requestPayload struct {
	Method    string                 `json:"method"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int                    `json:"tag,omitempty"`
}

func (c *Controller) request(method string, arguments map[string]interface{}) (err error) {
	// Prepare encoding pipeline
	tag := c.rnd.Int()
	pOut, pIn := io.Pipe()
	enc := json.NewEncoder(pIn)
	// Prepare the request
	var req *http.Request
	if req, err = http.NewRequest("POST", c.url, pOut); err != nil {
		err = fmt.Errorf("can't prepare request for '%s' method: %v", method, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(csrfHeader, c.sessionID)
	// Prepare encoding goroutine
	var encErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		encErr = enc.Encode(&requestPayload{
			Method:    method,
			Arguments: arguments,
			Tag:       tag,
		})
		pIn.Close()
		fmt.Println("DEBUG encoding goroutine done")
		wg.Done()
	}()
	// Execute request
	var resp *http.Response
	if resp, err = c.httpC.Do(req); err != nil {
		err = fmt.Errorf("request error: %v", err)
		wg.Wait()
		if encErr != nil {
			err = fmt.Errorf("%s | encoding error: %s", err, encErr)
		}
		return
	}
	defer resp.Body.Close()
	// Is the CRSF token invalid ?
	if resp.StatusCode == http.StatusConflict {
		c.sessionID = resp.Header.Get(csrfHeader)
		fmt.Printf("DEBUG updated CSRF token to: %s\n", c.sessionID)
		return c.request(method, arguments)
	}
	// Is request successfull ?
	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed with http code %d", resp.StatusCode)
		return
	}
	// Decode body
	var answer interface{}
	if err = json.NewDecoder(resp.Body).Decode(answer); err != nil {
		err = fmt.Errorf("can't decode request answer body as JSON: %v", err)
		return
	}
	// tmp
	fmt.Println(answer)
	return
}
