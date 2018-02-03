package TransmissionRPC

import (
	"fmt"
	"reflect"
)

var validTorrentFields []string

func init() {
	val := reflect.ValueOf(Torrent{})
	for i := 0; i < val.Type().NumField(); i++ {
		validTorrentFields = append(validTorrentFields, val.Type().Field(i).Tag.Get("json"))
	}
}

// TorrentGet returns the given of fields (mandatory) for each ids (optionnal)
func (c *Controller) TorrentGet(fields []string, ids []int) (torrents []*Torrent, err error) {
	// Validate fields
	var fieldInvalid bool
	var knownField string
	for _, inputField := range fields {
		fieldInvalid = true
		for _, knownField = range validTorrentFields {
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
	// Prepare
	arguments := torrentGetParams{
		Fields: fields,
		IDs:    ids,
	}
	var result torrentResults
	// Execute
	if err = c.rpcCall("torrent-get", &arguments, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %v", err)
		return
	}
	torrents = result.Torrents
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int    `json:"ids,omitempty"`
}

type torrentResults struct {
	Torrents []*Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent
type Torrent struct {
	ActivityDate *int `json:"activityDate"`
	AddedDate    *int `json:"addedDate"`
}
