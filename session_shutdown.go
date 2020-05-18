package transmissionrpc

import (
	"fmt"
)

/*
	Session shutdown
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L593
*/

// SessionClose tells the transmission session to shut down.
func (c *Client) SessionClose() (err error) {
	// Send request
	if err = c.rpcCall("session-close", nil, nil); err != nil {
		err = fmt.Errorf("'session-close' rpc method failed: %v", err)
	}
	return
}
