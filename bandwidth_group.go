package transmissionrpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

/*
	Bandwidth Group Accessors
	https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#48-bandwidth-groups
*/

// BandwidthGroup represents all possible fields of data for a bandwidth group.
type BandwidthGroup struct {
	HonorSessionLimits    bool   `json:"honorSessionLimits"`
	Name                  string `json:"name"`
	SpeedLimitDownEnabled bool   `json:"speed-limit-down-enabled"`
	SpeedLimitDown        int64  `json:"speed-limit-down"`
	SpeedLimitUpEnabled   bool   `json:"speed-limit-up-enabled"`
	SpeedLimitUp          int64  `json:"speed-limit-up"`
}

// BandwidthGroupGet returns the given fields for each bandwidth group provided. If no group names are provided, returns all groups.
// As of Transmission 4.0.3 it seems that filtering with group names does not work as documented here:
// https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#482-bandwidth-group-accessor-group-get
// Testing always yield the full list of bandwidth groups, even when submitting a filter.
func (c *Client) BandwidthGroupGet(ctx context.Context, groups []string) (bandwidthGroups []BandwidthGroup, err error) {
	var (
		filter string
		answer bandwidthGroupGetAnswer
	)
	if len(groups) > 0 {
		filter = strings.Join(groups, ",")
	}
	if err = c.rpcCall(ctx, "group-get", &bandwidthGroupGetParams{
		Group: filter,
	}, &answer); err != nil {
		err = fmt.Errorf("'group-get' rpc method failed: %w", err)
		return
	}
	bandwidthGroups = answer.Group
	return
}

type bandwidthGroupGetParams struct {
	Group string `json:"group,omitempty"`
}

type bandwidthGroupGetAnswer struct {
	Group []BandwidthGroup `json:"group"`
}

// BandwidthGroupSet applies a list of mutator(s) to a bandwidth group.
func (c *Client) BandwidthGroupSet(ctx context.Context, bwGroup BandwidthGroup) (err error) {
	// Validate
	if bwGroup.Name == "" {
		return errors.New("Bandwidth group must have a name")
	}
	// Send payload
	if err = c.rpcCall(ctx, "group-set", bwGroup, nil); err != nil {
		err = fmt.Errorf("'group-set' rpc method failed: %w", err)
	}
	return
}
