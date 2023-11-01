package client

import (
	"github.com/dustin/go-humanize"
	"time"
)

var magnitudes = []humanize.RelTimeMagnitude{
	{time.Second, "now", time.Second},
	{2 * time.Second, "1s %s", 1},
	{time.Minute, "%ds %s", time.Second},
	{2 * time.Minute, "1m %s", 1},
	{time.Hour, "%dm %s", time.Minute},
	{2 * time.Hour, "1h %s", 1},
	{humanize.Day, "%dh %s", time.Hour},
	{2 * humanize.Day, "1d %s", 1},
	{humanize.Week, "%dd %s", humanize.Day},
}
