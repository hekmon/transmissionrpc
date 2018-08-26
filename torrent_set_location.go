package transmissionrpc

import (
	"errors"
	"fmt"
)

/*
	Moving a torrent
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L421
*/

// TorrentSetLocation allows to set a new location for one or more torrents.
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L423
func (c *Client) TorrentSetLocation(payload *TorrentSetLocationPayload) (err error) {
	// Validate
	if payload == nil {
		return errors.New("payload can't be nil")
	}
	if len(payload.IDs) == 0 {
		return errors.New("there must be at least one ID")
	}
	// Send payload
	if err = c.rpcCall("torrent-set-location", payload, nil); err != nil {
		err = fmt.Errorf("'torrent-set-location' rpc method failed: %v", err)
	}
	return
}

// TorrentSetLocationPayload describes the torrents' id(s) and other options.
// https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L427
type TorrentSetLocationPayload struct {
	IDs      []int64 `json:"ids"`      // torrent list
	Location string  `json:"location"` // the new torrent location
	Move     bool    `json:"move"`     // if true, move from previous location. Otherwise, search "location" for files
}
