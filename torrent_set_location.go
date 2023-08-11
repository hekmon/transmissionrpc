package transmissionrpc

import (
	"context"
	"fmt"
)

/*
	Moving a torrent
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#36-moving-a-torrent
*/

// TorrentSetLocation allows to set a new location for one or more torrents.
// 'location' is the new torrent location.
// 'move' if true, move from previous location. Otherwise, search "location" for file.
func (c *Client) TorrentSetLocation(ctx context.Context, id int64, location string, move bool) (err error) {
	if err = c.rpcCall(ctx, "torrent-set-location", torrentSetLocationPayload{
		IDs:      []int64{id},
		Location: location,
		Move:     move,
	}, nil); err != nil {
		err = fmt.Errorf("'torrent-set-location' rpc method failed: %w", err)
	}
	return
}

// TorrentSetLocationHash allows to set a new location for one or more torrents.
// 'location' is the new torrent location.
// 'move' if true, move from previous location. Otherwise, search "location" for file.
func (c *Client) TorrentSetLocationHash(ctx context.Context, hash, location string, move bool) (err error) {
	if err = c.rpcCall(ctx, "torrent-set-location", torrentSetLocationHashPayload{
		Hashes:   []string{hash},
		Location: location,
		Move:     move,
	}, nil); err != nil {
		err = fmt.Errorf("'torrent-set-location' rpc method failed: %w", err)
	}
	return
}

type torrentSetLocationPayload struct {
	IDs      []int64 `json:"ids"`      // torrent list
	Location string  `json:"location"` // the new torrent location
	Move     bool    `json:"move"`     // if true, move from previous location. Otherwise, search "location" for files
}

type torrentSetLocationHashPayload struct {
	Hashes   []string `json:"ids"`      // torrent list
	Location string   `json:"location"` // the new torrent location
	Move     bool     `json:"move"`     // if true, move from previous location. Otherwise, search "location" for files
}
