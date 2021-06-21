package espn

import (
	"testing"
)

func TestExtractMatchIDFromLink(t *testing.T) {
	_ = []struct {
		m    ESPNMatchDescription
		name string
	}{
		{
			m: ESPNMatchDescription{
				Link: "/series/psl-2020-21-2021-1238103/islamabad-united-vs-multan-sultans-qualifier-1247041/live-cricket-score",
			},
			name: "normal",
		},
	}

	// m.extractMatchIDFromLink()

}
