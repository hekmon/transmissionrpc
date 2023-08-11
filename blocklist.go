package transmissionrpc

import (
	"context"
	"fmt"
)

/*
	Blocklist
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#43-blocklist
*/

// BlocklistUpdate triggers a blocklist update. It returns the number of entries of the updated blocklist.
func (c *Client) BlocklistUpdate(ctx context.Context) (nbEntries int64, err error) {
	var answer blocklistUpdateAnswer
	// Send request
	if err = c.rpcCall(ctx, "blocklist-update", nil, &answer); err == nil {
		nbEntries = answer.NbEntries
	} else {
		err = fmt.Errorf("'blocklist-update' rpc method failed: %w", err)
	}
	return
}

type blocklistUpdateAnswer struct {
	NbEntries int64 `json:"blocklist-size"`
}
