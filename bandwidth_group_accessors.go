package transmissionrpc

import (
	"context"
	"fmt"
	"reflect"
)

/*
   		Bandwidth Group Accessors
        https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md#482-bandwidth-group-accessor-group-get
*/

var validBandwidthGroupFields []string

func init() {
	bandwidthGroupType := reflect.TypeOf(BandwidthGroup{})
	for i := 0; i < bandwidthGroupType.NumField(); i++ {
		validBandwidthGroupFields = append(validBandwidthGroupFields, bandwidthGroupType.Field(i).Tag.Get("json"))
	}
}

// BandwidthGroupGetAll returns all the known fields for all bandwidth groups.
func (c *Client) BandwidthGroupGetAll(ctx context.Context) (bandwidthGroups []BandwidthGroup, err error) {
	// Send already validated fields to the low level function
	return c.bandwidthGroupGet(ctx, validBandwidthGroupFields, nil)
}

// BandwidthGroupGet returns the given fields for each bandwidth group provided. If no bandwidth group name is provided, returns
// the given fields for all bandwidth groups.
func (c *Client) BandwidthGroupGet(ctx context.Context, fields []string, groups []string) (bandwidthGroups []BandwidthGroup, err error) {
	if err = c.validateBandwidthGroupFields(fields); err != nil {
		return
	}
	return c.bandwidthGroupGet(ctx, validBandwidthGroupFields, groups)
}

// Validate fields
func (c *Client) validateBandwidthGroupFields(fields []string) (err error) {
	var fieldInvalid bool
	var knownField string
	for _, inputField := range fields {
		fieldInvalid = true
		for _, knownField = range validBandwidthGroupFields {
			if inputField == knownField {
				fieldInvalid = false
				break
			}
		}
		if fieldInvalid {
			err = fmt.Errorf("field '%s' is invalid", inputField)
			return
		}
	}
	return
}

// Low-level function for group-get RPC method call.
func (c *Client) bandwidthGroupGet(ctx context.Context, fields []string, groups []string) (bandwidthGroups []BandwidthGroup, err error) {
	var result bandwidthGroupGetResults
	if err = c.rpcCall(ctx, "group-get", &bandwidthGroupGetParams{
		Fields: fields,
		Names:  groups,
	}, &result); err != nil {
		err = fmt.Errorf("'group-get' rpc method failed: %w", err)
		return
	}
	bandwidthGroups = result.BandwidthGroups
	return
}

type bandwidthGroupGetParams struct {
	Fields []string `json:"fields,omitempty"`
	Names  []string `json:"names,omitempty"`
}

type bandwidthGroupGetResults struct {
	BandwidthGroups []BandwidthGroup `json:"group"`
}

// BandwidthGroup represents all possible fields of data for a bandwidth group.
type BandwidthGroup struct {
	HonorSessionLimits    bool   `json:"honorSessionLimits"`
	Name                  string `json:"name"`
	SpeedLimitDownEnabled bool   `json:"speed-limit-down-enabled"`
	SpeedLimitDown        int64  `json:"speed-limit-down"`
	SpeedLimitUpEnabled   bool   `json:"speed-limit-up-enabled"`
	SpeedLimitUp          int64  `json:"speed-limit-up"`
}
