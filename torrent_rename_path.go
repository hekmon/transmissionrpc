package transmissionrpc

import (
	"context"
	"fmt"
)

/*
	Rename a torrent path
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#37-renaming-a-torrents-path
*/

// TorrentRenamePath allows to rename torrent name or path.
// 'path' is the path to the file or folder that will be renamed.
// 'name' the file or folder's new name
func (c *Client) TorrentRenamePath(ctx context.Context, id int64, path, name string) (err error) {
	if err = c.rpcCall(ctx, "torrent-rename-path", torrentRenamePathPayload{
		IDs:  []int64{id},
		Path: path,
		Name: name,
	}, nil); err != nil {
		err = fmt.Errorf("'torrent-rename-path' rpc method failed: %w", err)
	}
	return
}

// TorrentRenamePathHash allows to rename torrent name or path by its hash.
func (c *Client) TorrentRenamePathHash(ctx context.Context, hash, path, name string) (err error) {
	if err = c.rpcCall(ctx, "torrent-rename-path", torrentRenamePathHashPayload{
		Hashes: []string{hash},
		Path:   path,
		Name:   name,
	}, nil); err != nil {
		err = fmt.Errorf("'torrent-rename-path' rpc method failed: %w", err)
	}
	return
}

type torrentRenamePathPayload struct {
	IDs  []int64 `json:"ids"`  // the torrent torrent list, as described in 3.1 (must only be 1 torrent)
	Path string  `json:"path"` // the path to the file or folder that will be renamed
	Name string  `json:"name"` // the file or folder's new name
}

type torrentRenamePathHashPayload struct {
	Hashes []string `json:"ids"`  // the torrent torrent list, as described in 3.1 (must only be 1 torrent)
	Path   string   `json:"path"` // the path to the file or folder that will be renamed
	Name   string   `json:"name"` // the file or folder's new name
}
