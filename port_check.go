package transmissionrpc

import (
	"fmt"
)

/*
	Port Checking
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L584
*/

// CheckPort allows tests to see if your incoming peer port is accessible from the outside world.
func (c *Client) CheckPort() (open bool, err error) {
	var result transmissionCheckPortAnswer
	// Send request
	if err = c.rpcCall("port-test", nil, &result); err == nil {
		open = result.PortOpen
	} else {
		err = fmt.Errorf("'port-test' rpc method failed: %v", err)
	}
	return
}

type transmissionCheckPortAnswer struct {
	PortOpen bool `json:"port-is-open"`
}
