package transmissionrpc

import (
	"errors"
	"fmt"
)

/*
	Rename a torrent path
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L438
*/

// TorrentRenamePath allows to rename torrent name or path.
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L438
func (c *Client) TorrentRenamePath(payload *TorrentRenamePathPayload) (err error) {
	// Validate
	if payload == nil {
		return errors.New("payload can't be nil")
	}
	if len(payload.IDs) == 0 {
		return errors.New("there must be at least one ID")
	}
	// Send payload
	if err = c.rpcCall("torrent-rename-path", payload, nil); err != nil {
		err = fmt.Errorf("'torrent-rename-path' rpc method failed: %v", err)
	}
	return
}

// TorrentSetLocationPayload describes the torrents' id(s) and other options.
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L427
type TorrentRenamePathPayload struct {
	IDs      []int64	`json:"ids"`	// torrent list
	Path 	 string 	`json:"path"`	// the new torrent path
	Name     string		`json:"name"`	// new torrent name
}
