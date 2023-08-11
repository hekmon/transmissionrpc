package transmissionrpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

/*
   Bandwidth Group Mutators
   https://github.com/transmission/transmission/blob/4.0.2/docs/rpc-spec.md#481-bandwidth-group-mutator-group-set
*/

// BandwidthGroupSet applies a list of mutator(s) to a bandwidth group.
func (c *Client) BandwidthGroupSet(ctx context.Context, payload BandwidthGroupSetPayload) (err error) {
	// Validate
	if reflect.ValueOf(payload.Name).Kind() == reflect.String {
		return errors.New("there must be one bandwidth group name")
	}
	// Send payload
	if err = c.rpcCall(ctx, "group-set", payload, nil); err != nil {
		err = fmt.Errorf("'group-set' rpc method failed: %w", err)
	}
	return
}

// BandwidthGroupSetPayload contains all the mutators applicable on one bandwidth group.
type BandwidthGroupSetPayload struct {
	HonorSessionLimits    *bool   `json:"honorSessionLimits"`       // true if session upload limits are honored
	Name                  *string `json:"name"`                     // Bandwidth group name
	SpeedLimitDownEnabled *bool   `json:"speed-limit-down-enabled"` // true means enabled
	SpeedLimitDown        *int64  `json:"speed-limit-down"`         // max global download speed (KBps)
	SpeedLimitUpEnabled   *bool   `json:"speed-limit-up-enabled"`   // true means enabled
	SpeedLimitUp          *int64  `json:"speed-limit-up"`           // max global upload speed (KBps)
}
