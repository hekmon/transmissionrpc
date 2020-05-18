package transmissionrpc

import (
	"fmt"
)

/*
	Queue Movement Requests
	https://github.com/transmission/transmission/blob/2.9x/extras/rpc-spec.txt#L601
*/

// QueueMoveTop moves IDs to the top of the queue list.
func (c *Client) QueueMoveTop(IDs []int64) (err error) {
	payload := &queueMovePayload{IDs: IDs}
	if err = c.rpcCall("queue-move-top", payload, nil); err != nil {
		err = fmt.Errorf("'queue-move-top' rpc method failed: %v", err)
	}
	return
}

// QueueMoveUp moves IDs of one position up on the queue list.
func (c *Client) QueueMoveUp(IDs []int64) (err error) {
	payload := &queueMovePayload{IDs: IDs}
	if err = c.rpcCall("queue-move-up", payload, nil); err != nil {
		err = fmt.Errorf("'queue-move-up' rpc method failed: %v", err)
	}
	return
}

// QueueMoveDown moves IDs of one position down on the queue list.
func (c *Client) QueueMoveDown(IDs []int64) (err error) {
	payload := &queueMovePayload{IDs: IDs}
	if err = c.rpcCall("queue-move-down", payload, nil); err != nil {
		err = fmt.Errorf("'queue-move-down' rpc method failed: %v", err)
	}
	return
}

// QueueMoveBottom moves IDs to the bottom of the queue list.
func (c *Client) QueueMoveBottom(IDs []int64) (err error) {
	payload := &queueMovePayload{IDs: IDs}
	if err = c.rpcCall("queue-move-bottom", payload, nil); err != nil {
		err = fmt.Errorf("'queue-move-bottom' rpc method failed: %v", err)
	}
	return
}

type queueMovePayload struct {
	IDs []int64 `json:"ids"`
}
