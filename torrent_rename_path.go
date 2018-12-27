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
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L440
func (c *Client) TorrentRenamePath(payload *TorrentRenamePathPayload) (err error) {
	// Validate
	if payload == nil {
		return errors.New("payload can't be nil")
	}
	if len(payload.IDs) != 1 {
		return errors.New("there must be one and only one ID")
	}
	// Send payload
	if err = c.rpcCall("torrent-rename-path", payload, nil); err != nil {
		err = fmt.Errorf("'torrent-rename-path' rpc method failed: %v", err)
	}
	return
}

// TorrentRenamePathPayload describes the torrents' id(s) and other options.
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L447
type TorrentRenamePathPayload struct {
	IDs  []int64 `json:"ids"`  // the torrent torrent list, as described in 3.1 (must only be 1 torrent)
	Path string  `json:"path"` // the path to the file or folder that will be renamed
	Name string  `json:"name"` // the file or folder's new name
}
