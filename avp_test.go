package avp

import "testing"

var t1Formats = Formats{
	{Resolution: 1080, AudioBitrate: 196},
	{Resolution: 4320},
	{Resolution: 2160},
	{Resolution: 720},
	{AudioBitrate: 48},
	{AudioBitrate: 316},
	{AudioBitrate: 128},
}

var t1Result = map[Quality]Formats{
	1: {
		{AudioBitrate: 316},
		{Resolution: 4320},
	},
	2: {
		{Resolution: 1080, AudioBitrate: 196},
	},
	3: {
		{AudioBitrate: 128},
		{Resolution: 720},
	},
	4: {
		{AudioBitrate: 128},
		{Resolution: 720},
	},
}

func TestAll(t *testing.T) {
	avp := New(t1Formats)
	check(t, avp.qltMap, t1Result)
}

func check(t *testing.T, got, res map[Quality]Formats) {
	for i, g := range got {
		equal(t, g, res[i])
	}
}

func equal(t *testing.T, a, b Formats) bool {
	if len(a) != len(b) {
		return false
	}

	for i, aa := range a {
		if !fmtEqual(aa, b[i]) {
			t.Error("Expected:", b[i], "Got:", aa)
			return false
		}
	}

	return true
}

func fmtEqual(a, b *Format) bool {
	// video
	if a.Resolution != b.Resolution || a.VideoCodec != b.VideoCodec || a.VideoBitrate != b.VideoBitrate || a.VideoHDR != b.VideoHDR || a.VideoHFR != b.VideoHDR {
		return false
	}

	if a.AudioBitrate != b.AudioBitrate || a.AudioChannels != b.AudioChannels || a.AudioCodec != b.AudioCodec || a.AudioVBR != b.AudioVBR {
		return false
	}
	return true
}
