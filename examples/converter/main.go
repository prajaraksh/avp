package main

import (
	"fmt"

	"github.com/prajaraksh/avp"
)

type ifmt struct {
	URL          string
	Resolution   int
	AudioBitrate int
}

var ifmtData = []*ifmt{
	{"a", 1080, 196},
	{"b", 4320, 0},
	{"c", 0, 320},
	{"d", 2340, 231},
	{"e", 720, 128},
}

type ifmts []*ifmt

func (ifs ifmts) Len() int {
	return len(ifs)
}

func (ifs ifmts) ToFormat(id int) *avp.Format {
	return &avp.Format{ID: id, Resolution: ifs[id].Resolution, AudioBitrate: ifs[id].AudioBitrate}
}

func main() {
	cavp := avp.WithConverter(ifmts(ifmtData), avp.AVProfile)
	res := cavp.Best()
	fmt.Println(res)
}
