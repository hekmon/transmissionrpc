package transmissionrpc

import (
	"context"
	"fmt"

	"github.com/hekmon/cunits/v2"
)

/*
	Session Statistics
    https://github.com/transmission/transmission/blob/4.0.3/docs/rpc-spec.md#42-session-statistics
*/

// SessionStats returns all (current/cumulative) statistics.
func (c *Client) SessionStats(ctx context.Context) (stats SessionStats, err error) {
	if err = c.rpcCall(ctx, "session-stats", nil, &stats); err != nil {
		err = fmt.Errorf("'session-stats' rpc method failed: %w", err)
	}
	return
}

// SessionStats represents all (current/cumulative) statistics.
type SessionStats struct {
	ActiveTorrentCount int64               `json:"activeTorrentCount"`
	DownloadSpeed      int64               `json:"downloadSpeed"`
	PausedTorrentCount int64               `json:"pausedTorrentCount"`
	TorrentCount       int64               `json:"torrentCount"`
	UploadSpeed        int64               `json:"uploadSpeed"`
	CumulativeStats    SessionStatsDetails `json:"cumulative-stats"`
	CurrentStats       SessionStatsDetails `json:"current-stats"`
}

// CumulativeStats is subset of SessionStats.
type SessionStatsDetails struct {
	DownloadedBytes int64 `json:"downloadedBytes"`
	FilesAdded      int64 `json:"filesAdded"`
	SecondsActive   int64 `json:"secondsActive"`
	SessionCount    int64 `json:"sessionCount"`
	UploadedBytes   int64 `json:"uploadedBytes"`
}

// GetDownloaded returns cumulative stats downloaded size in a handy format
func (cs *SessionStatsDetails) GetDownloaded() (downloaded cunits.Bits) {
	return cunits.ImportInByte(float64(cs.DownloadedBytes))
}

// GetUploaded returns cumulative stats uploaded size in a handy format
func (cs *SessionStatsDetails) GetUploaded() (uploaded cunits.Bits) {
	return cunits.ImportInByte(float64(cs.UploadedBytes))
}
