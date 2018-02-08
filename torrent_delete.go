package transmissionrpc

import (
	"errors"
	"fmt"
)

/*
	Removing a Torrent
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L392
*/

// TorrentDelete apply a list of mutator(s) to a list of torrent ids.
func (c *Controller) TorrentDelete(payload *TorrentDeletePayload) (err error) {
	// Validate
	if payload == nil {
		err = errors.New("payload can't be nil")
		return
	}
	// Send payload
	if err = c.rpcCall("torrent-remove", payload, nil); err != nil {
		err = fmt.Errorf("'torrent-remove' rpc method failed: %v", err)
		return
	}
	return
}

// TorrentDeletePayload allows to delete several torrents at once.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L398
type TorrentDeletePayload struct {
	IDs             []int64 `json:"ids"`
	DeleteLocalData bool    `json:"delete-local-data"`
}
