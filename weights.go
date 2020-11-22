package avp

var codecWeight = map[string]int{

	// video
	"h264-baseline": 20, // h264
	"h264-main":     25,
	"h264":          30, // generic
	"h264-high":     40,
	"h265":          45,

	"vp8": 25, // vp9
	"vp9": 45,

	"av1": 50, // av1

	// audio
	"mp3": 20, // mp3

	"vorbis": 25, // vorbis

	"aac-he": 25, // aac
	"aac":    30, // generic
	"aac-lc": 35,

	"opus": 40, // opus
}

var channelWeight = map[int]int{
	2: 20,
	6: 30,
	8: 40,
}

// Profile to define selection standards
type Profile struct {
	High, Medium, Low *Format
}

// AVProfile ,Audio Video Profile
var AVProfile = Profile{
	High:   &Format{Resolution: 1080, AudioBitrate: 256},
	Medium: &Format{Resolution: 720, AudioBitrate: 192},
	Low:    &Format{Resolution: 360, AudioBitrate: 128},
}
