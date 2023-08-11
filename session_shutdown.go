package transmissionrpc

import (
	"context"
	"fmt"
)

/*
	Session shutdown
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#45-session-shutdown
*/

// SessionClose tells the transmission session to shut down.
func (c *Client) SessionClose(ctx context.Context) (err error) {
	// Send request
	if err = c.rpcCall(ctx, "session-close", nil, nil); err != nil {
		err = fmt.Errorf("'session-close' rpc method failed: %w", err)
	}
	return
}
