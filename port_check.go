package transmissionrpc

import (
	"fmt"
)

/*
	Port Checking
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L569
*/

// CheckPort allow to see if your incoming peer port is accessible from the outside world.
func (c *Client) CheckPort() (opened bool, err error) {
	var result transmissionCheckPortAnswer
	// Send request
	if err = c.rpcCall("port-test", nil, &result); err != nil {
		err = fmt.Errorf("'port-test' rpc method failed: %v", err)
		return
	}
	opened = result.PortOpened
	return
}

type transmissionCheckPortAnswer struct {
	PortOpened bool `json:"port-is-open"`
}
