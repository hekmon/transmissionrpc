package transmissionrpc

import (
	"fmt"
)

/*
	Free Space
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L603
*/

// FreeSpace allow to see how much free space is available in a client-specified folder.
func (c *Client) FreeSpace(path string) (freeBytes int64, err error) {
	payload := &transmissionFreeSpacePayload{Path: path}
	var space TransmissionFreeSpace
	if err = c.rpcCall("free-space", payload, &space); err != nil {
		err = fmt.Errorf("'free-space' rpc method failed: %v", err)
		return
	}
	if space.Path != path {
		err = fmt.Errorf("returned path '%s' does not match with requested path '%s'", space.Path, path)
		return
	}
	freeBytes = space.Size
	return
}

type transmissionFreeSpacePayload struct {
	Path string `json:"path"`
}

// TransmissionFreeSpace represents the freespace available in bytes for a specific path.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L616
type TransmissionFreeSpace struct {
	Path string `json:"path"`
	Size int64  `json:"size-bytes"`
}
