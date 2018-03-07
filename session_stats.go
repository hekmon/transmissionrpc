package transmissionrpc

import (
	"fmt"
)

/*
	Session Statistics
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L531
*/

// SessionStats returns all (current/cumulative) statistics.
func (c *Client) SessionStats() (stats *SessionStats, err error) {
	if err = c.rpcCall("session-stats", nil, &stats); err != nil {
		err = fmt.Errorf("'session-stats' rpc method failed: %v", err)
	}
	return
}

// SessionStats represents all (current/cumulative) statistics.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L539
type SessionStats struct {
	ActiveTorrentCount int64            `json:"activeTorrentCount"`
	CumulativeStats    *CumulativeStats `json:"cumulative-stats"`
	CurrentStats       *CurrentStats    `json:"current-stats"`
	DownloadSpeed      int64            `json:"downloadSpeed"`
	PausedTorrentCount int64            `json:"pausedTorrentCount"`
	TorrentCount       int64            `json:"torrentCount"`
	UploadSpeed        int64            `json:"uploadSpeed"`
}

// CumulativeStats is subset of SessionStats.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L547
type CumulativeStats struct {
	DownloadedBytes int64 `json:"downloadedBytes"`
	FilesAdded      int64 `json:"filesAdded"`
	SecondsActive   int64 `json:"secondsActive"`
	SessionCount    int64 `json:"sessionCount"`
	UploadedBytes   int64 `json:"uploadedBytes"`
}

// CurrentStats is subset of SessionStats.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L555
type CurrentStats struct {
	DownloadedBytes int64 `json:"downloadedBytes"`
	FilesAdded      int64 `json:"filesAdded"`
	SecondsActive   int64 `json:"secondsActive"`
	SessionCount    int64 `json:"sessionCount"`
	UploadedBytes   int64 `json:"uploadedBytes"`
}
