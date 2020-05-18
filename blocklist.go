package transmissionrpc

import (
	"fmt"
)

/*
	Blocklist
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L578
*/

// BlocklistUpdate triggers a blocklist update. It returns the number of entries of the updated blocklist.
func (c *Client) BlocklistUpdate() (nbEntries int64, err error) {
	var answer blocklistUpdateAnswer
	// Send request
	if err = c.rpcCall("blocklist-update", nil, &answer); err == nil {
		nbEntries = answer.nbEntries
	} else {
		err = fmt.Errorf("'blocklist-update' rpc method failed: %v", err)
	}
	return
}

type blocklistUpdateAnswer struct {
	nbEntries int64 `json:"blocklist-size"`
}
