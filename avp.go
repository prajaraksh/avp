package avp

import (
	"fmt"
	"math"
	"sort"
)

// Format of a video or audio file
type Format struct {
	ID         int
	Resolution int // Video
	// VideoBitrate in Kbps
	VideoBitrate int
	VideoCodec   string
	VideoHFR     bool
	VideoHDR     bool
	// AudioBitrate in Kbps
	AudioBitrate  int // Audio
	AudioChannels int
	AudioCodec    string
	AudioVBR      bool
}

// Formats is slice of Format pointers
type Formats []*Format

// AVP ,Audio Video Picker
type AVP struct {
	fs     Formats
	p      Profile
	qltMap map[Quality]Formats
}

// New , given formats are categorized into different sections
func New(fs Formats) *AVP {
	return WithProfile(fs, AVProfile)
}

// WithProfile takes Profile option
func WithProfile(fs Formats, p Profile) *AVP {
	avp := &AVP{
		fs:     fs,
		p:      p,
		qltMap: make(map[Quality]Formats, 5),
	}

	// segreate list based on their type
	avfs, aofs, vofs := segregateByType(fs)

	// sort each list
	sort.Sort(formats(avfs))
	sort.Sort(formats(aofs))
	sort.Sort(formats(vofs))

	avpfs := segregateByProfile(avfs, avp.p)
	aopfs := segregateByProfile(aofs, avp.p)
	vopfs := segregateByProfile(vofs, avp.p)

	finalSelect(avpfs, aopfs, vopfs, avp.qltMap)

	return avp
}

// Best ,Formats greater than `High` profile
func (avp *AVP) Best() Formats {
	return avp.qltMap[Best]
}

// High ,Formats less than or equal to `High` Profile and greater than `Medium` profile
func (avp *AVP) High() Formats {
	return avp.qltMap[High]
}

// Medium ,Formats less than or equal to `Medium` Profile and greater than `Low` profile
func (avp *AVP) Medium() Formats {
	return avp.qltMap[Medium]
}

// Low ,Formats less than or equal to `Low` profile
func (avp *AVP) Low() Formats {
	return avp.qltMap[Low]
}

// Converter should be implemented by every site
type Converter interface {
	Len() int
	ToFormat(id int) *Format
}

// WithConverter takes in `conv` which converters to internal format
// as required
func WithConverter(conv Converter, p Profile) *AVP {
	l := conv.Len()

	fs := make(Formats, l)

	for id := 0; id < l; id++ {
		fs[id] = conv.ToFormat(id)
		fs[id].ID = id
	}

	avp := WithProfile(fs, p)

	return avp
}

// Quality type represents type of profile
type Quality int

const (
	// None represents, no selection of particular quality
	None Quality = iota
	// Best represents best available quality
	Best
	// High ,high quality based on profile
	High
	// Medium ,medium quality based on profile
	Medium
	// Low ,low quality based on profile
	Low
)

// OfQuality is helper to give formats of particular quality
func (avp *AVP) OfQuality(qlty Quality) Formats {
	return avp.qltMap[qlty]
}

// Weight = Resolution*10 + VideoBitrate + 100*VideoCodec + 1000 + 1000
// Weight = AudioBitrate + ChannelWeight + AudioCodec + 40

func weight(f *Format) int {
	vw := f.Resolution*10 + f.VideoBitrate + 100*codecWeight[f.VideoCodec]
	if f.VideoHDR {
		vw += 1000
	}
	if f.VideoHFR {
		vw += 1000
	}

	aw := f.AudioBitrate + channelWeight[f.AudioChannels] + codecWeight[f.AudioCodec]

	if f.AudioVBR {
		aw += 40
	}

	return vw + aw
}

func finalSelect(avpfs, aopfs, vopfs []Formats, m map[Quality]Formats) {
	for i := 1; i < 5; i++ {
		m[Quality(i)] = selectForSection(i, avpfs, aopfs, vopfs)
	}
}

func selectForSection(id int, avpfs, aopfs, vopfs []Formats) Formats {
	avf, aof, vof := someFormat(id, avpfs), someFormat(id, aopfs), someFormat(id, vopfs)

	avgf, aogf, vogf := avf != nil, aof != nil, vof != nil

	if avgf && aogf && vogf {
		return best(avf, aof, vof)
	} else if avgf && ((!aogf && !vogf) || (!aogf && vogf) || (aogf && !vogf)) {
		return Formats{avf}
	} else if !avgf && aogf && vogf {
		return Formats{aof, vof}
	} else if !avgf && !aogf && vogf {
		return Formats{vof}
	} else if !avgf && aogf && !vogf {
		return Formats{aof}
	}

	return Formats{}
}

func best(av, ao, vo *Format) Formats {
	if weight(av) < weight(ao)+weight(vo) {
		return Formats{ao, vo}
	}
	return Formats{av}
}

func someFormat(id int, pfs []Formats) *Format {
	// 0, downward - ++
	// i.e., we move from best, high, medium, low
	// since we are moving downward, we use `bestof()`
	// i.e., user is expecting best, we didn't have any,
	// so we choose, next best from high, if not from medium, ...
	// 1, upward - --
	// i.e., user requested from medium,
	// if medium and low are empty
	// we need `leastof()` from high, ...

	if res := selectFormat(id, 0, pfs); res != nil {
		return res
	}

	return selectFormat(id, 1, pfs)
}

func selectFormat(id, motion int, pfs []Formats) *Format {
	var fn func(fs Formats) *Format
	var op func()

	if motion == 0 {
		op = func() { id++ }
		fn = bestof
	} else {
		op = func() { id-- }
		fn = leastof
	}

	for ; id > 0 && id < 5; op() {
		res := fn(pfs[id])
		if res != nil {
			return res
		}
	}

	return nil
}

func bestof(fs Formats) *Format {
	if len(fs) != 0 {
		// last format is the best quality
		return fs[len(fs)-1]
	}
	return nil
}

func leastof(fs Formats) *Format {
	if len(fs) != 0 {
		// first format is of least quality
		return fs[0]
	}
	return nil
}

// audio & video format list
// audio only format list
// video only format list
func segregateByType(fs Formats) (avfs Formats, aofs Formats, vofs Formats) {

	for _, f := range fs {

		var ao, vo bool
		ao = audioOnly(f)
		vo = videoOnly(f)

		if ao && vo {
			avfs = append(avfs, f)
		} else if ao {
			aofs = append(aofs, f)
		} else if vo {
			vofs = append(vofs, f)
		}

	}

	return
}

func audioOnly(f *Format) bool {
	if f.AudioBitrate != 0 || f.AudioChannels != 0 || f.AudioCodec != "" || f.AudioVBR {
		return true
	}
	return false
}

func videoOnly(f *Format) bool {
	if f.Resolution != 0 || f.VideoBitrate != 0 || f.VideoCodec != "" || f.VideoHFR || f.VideoHDR {
		return true
	}
	return false
}

func segregateByProfile(fs Formats, p Profile) []Formats {

	// 0 - None
	// 1 - best
	// 2 - high
	// 3 = medium
	// 4 = low
	profileMatchers := make([]Formats, 5)

	bestProfile := &Format{Resolution: math.MaxInt64, AudioBitrate: math.MaxInt64}
	profileMatchers[1] = matches(bestProfile, p.High, fs)
	profileMatchers[2] = matches(p.High, p.Medium, fs)
	profileMatchers[3] = matches(p.Medium, p.Low, fs)
	profileMatchers[4] = matches(p.Low, &Format{}, fs)

	return profileMatchers
}

func matches(curProfile, nextProfile *Format, fs Formats) Formats {

	matchedFormats := make(Formats, 0)

	for _, f := range fs {
		if res := match(curProfile, nextProfile, f); res != nil {
			matchedFormats = append(matchedFormats, res)
		}
	}

	return matchedFormats
}

// match
// Resolution    - <=
// VideoBitrate  - <=
// VideoCodec    - ==
// VideoHFR      - ==
// VideoHDR      - ==
// AudioBitrate  - <=
// AudioChannels - <=
// AudioCodec    - ==
// AudioVBR      - ==
func match(curProfile, nextProfile, f *Format) *Format {

	if curProfile.VideoCodec != "" {
		if curProfile.VideoCodec != f.VideoCodec {
			return nil
		}
	}

	if curProfile.Resolution != 0 && f.Resolution != 0 {
		if f.Resolution > curProfile.Resolution || f.Resolution < nextProfile.Resolution {
			return nil
		}
	}

	if curProfile.VideoBitrate != 0 && f.VideoBitrate != 0 {
		if f.VideoBitrate > curProfile.VideoBitrate {
			return nil
		}
	}

	if curProfile.VideoHFR {
		if !f.VideoHFR {
			return nil
		}
	}

	if curProfile.VideoHDR {
		if !f.VideoHDR {
			return nil
		}
	}

	if curProfile.AudioBitrate != 0 && f.AudioBitrate != 0 {
		if f.AudioBitrate > curProfile.AudioBitrate || f.AudioBitrate < nextProfile.AudioBitrate {
			return nil
		}
	}

	if curProfile.AudioChannels != 0 && f.AudioChannels != 0 {
		if f.AudioChannels != curProfile.AudioChannels {
			return nil
		}
	}

	if curProfile.AudioCodec != "" {
		if curProfile.AudioCodec != f.AudioCodec {
			return nil
		}
	}

	if curProfile.AudioVBR {
		if !f.AudioVBR {
			return nil
		}
	}

	return f
}

func (f *Format) String() string {
	ao, vo := audioOnly(f), videoOnly(f)

	if ao && vo {
		return fmt.Sprintf("ID:%d, Resolution:%dp, AudioBitrate:%dKbps\n", f.ID, f.Resolution, f.AudioBitrate)
	} else if ao {
		return fmt.Sprintf("ID:%d, AudioBitrate:%dKbps\n", f.ID, f.AudioBitrate)
	} else {
		return fmt.Sprintf("ID:%d, Resolution:%dp\n", f.ID, f.Resolution)
	}
}

type formats Formats

func (fs formats) Less(i, j int) bool {
	return weight(fs[i]) < weight(fs[j])
}

func (fs formats) Len() int {
	return len(fs)
}

func (fs formats) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs formats) String() string {
	str := ""
	for _, f := range fs {
		str += fmt.Sprintf("%v\n", f)
	}
	return str
}
