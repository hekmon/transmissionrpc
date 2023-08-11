package transmissionrpc

import (
	"context"
	"fmt"

	"github.com/hekmon/cunits/v2"
)

/*
	Free Space
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#47-free-space
*/

// FreeSpace allow to see how much free space is available in a client-specified folder.
func (c *Client) FreeSpace(ctx context.Context, path string) (freeSpace, totalSize cunits.Bits, err error) {
	payload := &transmissionFreeSpacePayload{Path: path}
	var space TransmissionFreeSpace
	if err = c.rpcCall(ctx, "free-space", payload, &space); err == nil {
		if space.Path == path {
			freeSpace = cunits.ImportInByte(float64(space.Size))
			totalSize = cunits.ImportInByte(float64(space.TotalSize))
		} else {
			err = fmt.Errorf("returned path '%s' does not match with requested path '%s'", space.Path, path)
		}
	} else {
		err = fmt.Errorf("'free-space' rpc method failed: %w", err)
	}
	return
}

type transmissionFreeSpacePayload struct {
	Path string `json:"path"`
}

// TransmissionFreeSpace represents the freespace available in bytes for a specific path.
type TransmissionFreeSpace struct {
	Path      string `json:"path"`
	Size      int64  `json:"size-bytes"`
	TotalSize int64  `json:"total_size"` // RPC v17
}
