package TransmissionRPC

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

/*
	Torrent Mutators
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L90
*/

// TorrentSet apply a list of mutator(s) to a list of torrent ids
func (c *Controller) TorrentSet(payload *TorrentSetPayload) (err error) {
	// Validate
	if payload == nil {
		return errors.New("payload can't be nil")
	}
	if len(payload.IDs) == 0 {
		return errors.New("there must be at least one ID")
	}
	// Send payload
	if err = c.rpcCall("torrent-set", payload, nil); err != nil {
		err = fmt.Errorf("'torrent-set' rpc method failed: %v", err)
	}
	return
}

// TorrentSetPayload contains all the mutators appliable on one torrent
type TorrentSetPayload struct {
	BandwidthPriority   *int64         `json:"bandwidthPriority"`
	DownloadLimit       *int64         `json:"downloadLimit"`
	DownloadLimited     *bool          `json:"downloadLimited"`
	FilesWanted         []int64        `json:"files-wanted"`   // empty array == all files
	FilesUnwanted       []int64        `json:"files-unwanted"` // empty array == all files
	HonorsSessionLimits *bool          `json:"honorsSessionLimits"`
	IDs                 []int64        `json:"ids"`
	Location            *string        `json:"location"`
	PeeLlimit           *int64         `json:"peer-limit"`
	PriorityHigh        []int64        `json:"priority-high"`   // empty array == all files
	PriorityLow         []int64        `json:"priority-low"`    // empty array == all files
	PriorityNormal      []int64        `json:"priority-normal"` // empty array == all files
	QueuePosition       *int64         `json:"queuePosition"`
	SeedIdleLimit       *time.Duration `json:"seedIdleLimit"` // unit is seconds
	SeedIdleMode        *int64         `json:"seedIdleMode"`
	SeedRatioLimit      *float64       `json:"seedRatioLimit"`
	SeedRatioMode       *int64         `json:"seedRatioMode"`
	TrackerAdd          []string       `json:"trackerAdd"`
	TrackerRemove       []int64        `json:"trackerRemove"`
	TrackerReplace      []string       `json:"trackerReplace"` // mmmm...
	UploadLimit         *int64         `json:"uploadLimit"`    // KBps
	UploadLimited       *bool          `json:"uploadLimited"`
}

// MarshalJSON allows to marshall into JSON only the non nil fields.
// It differs from 'omitempty' which also skip default values (as 0 which can valid here).
func (tsp *TorrentSetPayload) MarshalJSON() (data []byte, err error) {
	// Build an intermediary payload with base types
	type baseTorrentSetPayload TorrentSetPayload
	tmp := struct {
		SeedIdleLimit *int64 `json:"seedIdleLimit"`
		*baseTorrentSetPayload
	}{
		baseTorrentSetPayload: (*baseTorrentSetPayload)(tsp),
	}
	if tsp.SeedIdleLimit != nil {
		sil := int64(*tsp.SeedIdleLimit / time.Second)
		tmp.SeedIdleLimit = &sil
	}
	// Build a payload with only the non nil fields
	tspv := reflect.ValueOf(tmp)
	tspt := tspv.Type()
	cleanPayload := make(map[string]interface{}, tspt.NumField())
	var currentValue, nestedStruct, currentNestedValue reflect.Value
	var currentStructField, currentNestedStructField reflect.StructField
	var j int
	for i := 0; i < tspv.NumField(); i++ {
		currentValue = tspv.Field(i)
		currentStructField = tspt.Field(i)
		// For the nested/inherited struct
		if currentStructField.Name == "baseTorrentSetPayload" && !currentValue.IsNil() {
			nestedStruct = reflect.Indirect(currentValue)
			for j = 0; j < nestedStruct.NumField(); j++ {
				currentNestedValue = nestedStruct.Field(j)
				currentNestedStructField = nestedStruct.Type().Field(j)
				if !currentNestedValue.IsNil() {
					cleanPayload[currentNestedStructField.Tag.Get("json")] = currentNestedValue.Interface()
				}
			}
			continue
		}
		// For the overloaded fields
		if !currentValue.IsNil() {
			fmt.Println(currentStructField.Name)
			cleanPayload[currentStructField.Tag.Get("json")] = currentValue.Interface()
		}
	}
	// Marshall the clean payload
	return json.Marshal(cleanPayload)
}
