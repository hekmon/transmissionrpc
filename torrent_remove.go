package transmissionrpc

import (
	"context"
	"errors"
	"fmt"
)

/*
	Removing a Torrent
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#35-removing-a-torrent
*/

// TorrentRemove allows to delete one or more torrents only or with their data.
func (c *Client) TorrentRemove(ctx context.Context, payload TorrentRemovePayload) (err error) {
	// Validate
	if len(payload.IDs) == 0 {
		return errors.New("there must be at least one ID")
	}
	// Send payload
	if err = c.rpcCall(ctx, "torrent-remove", payload, nil); err != nil {
		return fmt.Errorf("'torrent-remove' rpc method failed: %w", err)
	}
	return
}

// TorrentRemovePayload holds the torrent id(s) to delete with a data deletion flag.
type TorrentRemovePayload struct {
	IDs             []int64 `json:"ids"`
	DeleteLocalData bool    `json:"delete-local-data"`
}
