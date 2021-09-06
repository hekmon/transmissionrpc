package transmissionrpc

import (
	"context"
	"fmt"
)

/*
	Torrent Action Requests
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L86
*/

type torrentActionIDsParam struct {
	IDs []int64 `json:"ids,omitempty"`
}

type torrentActionHashesParam struct {
	IDs []string `json:"ids,omitempty"`
}

type torrentActionRecentlyActiveParam struct {
	IDs string `json:"ids"`
}

// TorrentStartIDs starts torrent(s) which id is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStartIDs(ctx context.Context, ids []int64) (err error) {
	if err = c.rpcCall(ctx, "torrent-start", &torrentActionIDsParam{IDs: ids}, nil); err != nil {
		err = fmt.Errorf("'torrent-start' rpc method failed: %v", err)
	}
	return
}

// TorrentStartHashes starts torrent(s) which hash is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStartHashes(ctx context.Context, hashes []string) (err error) {
	if err = c.rpcCall(ctx, "torrent-start", &torrentActionHashesParam{IDs: hashes}, nil); err != nil {
		err = fmt.Errorf("'torrent-start' rpc method failed: %v", err)
	}
	return
}

// TorrentStartRecentlyActive starts torrent(s) which have been recently active.
func (c *Client) TorrentStartRecentlyActive(ctx context.Context) (err error) {
	if err = c.rpcCall(ctx, "torrent-start", &torrentActionRecentlyActiveParam{IDs: "recently-active"}, nil); err != nil {
		err = fmt.Errorf("'torrent-start' rpc method failed: %v", err)
	}
	return
}

// TorrentStartNowIDs starts (now) torrent(s) which id is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStartNowIDs(ctx context.Context, ids []int64) (err error) {
	if err = c.rpcCall(ctx, "torrent-start-now", &torrentActionIDsParam{IDs: ids}, nil); err != nil {
		err = fmt.Errorf("'torrent-start-now' rpc method failed: %v", err)
	}
	return
}

// TorrentStartNowHashes starts (now) torrent(s) which hash is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStartNowHashes(ctx context.Context, hashes []string) (err error) {
	if err = c.rpcCall(ctx, "torrent-start-now", &torrentActionHashesParam{IDs: hashes}, nil); err != nil {
		err = fmt.Errorf("'torrent-start-now' rpc method failed: %v", err)
	}
	return
}

// TorrentStartNowRecentlyActive starts (now) torrent(s) which have been recently active.
func (c *Client) TorrentStartNowRecentlyActive(ctx context.Context) (err error) {
	if err = c.rpcCall(ctx, "torrent-start-now", &torrentActionRecentlyActiveParam{IDs: "recently-active"}, nil); err != nil {
		err = fmt.Errorf("'torrent-start-now' rpc method failed: %v", err)
	}
	return
}

// TorrentStopIDs stops torrent(s) which id is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStopIDs(ctx context.Context, ids []int64) (err error) {
	if err = c.rpcCall(ctx, "torrent-stop", &torrentActionIDsParam{IDs: ids}, nil); err != nil {
		err = fmt.Errorf("'torrent-stop' rpc method failed: %v", err)
	}
	return
}

// TorrentStopHashes stops torrent(s) which hash is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentStopHashes(ctx context.Context, hashes []string) (err error) {
	if err = c.rpcCall(ctx, "torrent-stop", &torrentActionHashesParam{IDs: hashes}, nil); err != nil {
		err = fmt.Errorf("'torrent-stop' rpc method failed: %v", err)
	}
	return
}

// TorrentStopRecentlyActive stops torrent(s) which have been recently active.
func (c *Client) TorrentStopRecentlyActive(ctx context.Context) (err error) {
	if err = c.rpcCall(ctx, "torrent-stop", &torrentActionRecentlyActiveParam{IDs: "recently-active"}, nil); err != nil {
		err = fmt.Errorf("'torrent-stop' rpc method failed: %v", err)
	}
	return
}

// TorrentVerifyIDs verifys torrent(s) which id is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentVerifyIDs(ctx context.Context, ids []int64) (err error) {
	if err = c.rpcCall(ctx, "torrent-verify", &torrentActionIDsParam{IDs: ids}, nil); err != nil {
		err = fmt.Errorf("'torrent-verify' rpc method failed: %v", err)
	}
	return
}

// TorrentVerifyHashes verifys torrent(s) which hash is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentVerifyHashes(ctx context.Context, hashes []string) (err error) {
	if err = c.rpcCall(ctx, "torrent-verify", &torrentActionHashesParam{IDs: hashes}, nil); err != nil {
		err = fmt.Errorf("'torrent-verify' rpc method failed: %v", err)
	}
	return
}

// TorrentVerifyRecentlyActive verifys torrent(s) which have been recently active.
func (c *Client) TorrentVerifyRecentlyActive(ctx context.Context) (err error) {
	if err = c.rpcCall(ctx, "torrent-verify", &torrentActionRecentlyActiveParam{IDs: "recently-active"}, nil); err != nil {
		err = fmt.Errorf("'torrent-verify' rpc method failed: %v", err)
	}
	return
}

// TorrentReannounceIDs reannounces torrent(s) which id is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentReannounceIDs(ctx context.Context, ids []int64) (err error) {
	if err = c.rpcCall(ctx, "torrent-reannounce", &torrentActionIDsParam{IDs: ids}, nil); err != nil {
		err = fmt.Errorf("'torrent-reannounce' rpc method failed: %v", err)
	}
	return
}

// TorrentReannounceHashes reannounces torrent(s) which hash is in the provided slice.
// Can be one, can be several, can be all (if slice is empty or nil).
func (c *Client) TorrentReannounceHashes(ctx context.Context, hashes []string) (err error) {
	if err = c.rpcCall(ctx, "torrent-reannounce", &torrentActionHashesParam{IDs: hashes}, nil); err != nil {
		err = fmt.Errorf("'torrent-reannounce' rpc method failed: %v", err)
	}
	return
}

// TorrentReannounceRecentlyActive reannounces torrent(s) which have been recently active.
func (c *Client) TorrentReannounceRecentlyActive(ctx context.Context) (err error) {
	if err = c.rpcCall(ctx, "torrent-reannounce", &torrentActionRecentlyActiveParam{IDs: "recently-active"}, nil); err != nil {
		err = fmt.Errorf("'torrent-reannounce' rpc method failed: %v", err)
	}
	return
}
