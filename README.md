# AVP

Audio Video Picker, tries to choose best format available of given formats.

## Usage

```go
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

// internal format to avp format
func itof(ifmtData []*ifmt) []*avp.Format {
    t := make([]*avp.Format, 0, len(ifmtData))

    for id, f := range ifmtData {
        t = append(t, &avp.Format{ID: id, Resolution: f.Resolution, AudioBitrate: f.AudioBitrate})
    }
    return t
}

func main() {
    a := avp.New(itof(ifmtData))
    res := a.Best()
    fmt.Println(res)
}
```

```go
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

```
