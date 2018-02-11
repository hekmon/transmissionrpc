package transmissionrpc

import (
	"errors"
	"fmt"
)

/*
	Moving a torrent
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L406
*/

// TorrentSetLocation allows to set a new location for one or more torrents.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L408
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
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L412
type TorrentSetLocationPayload struct {
	IDs      []int64 `json:"ids"`      // torrent list
	Location string  `json:"location"` // the new torrent location
	Move     bool    `json:"move"`     // if true, move from previous location. Otherwise, search "location" for files
}
