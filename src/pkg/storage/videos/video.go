package videos

import (
	"bytes"
	"io"

	"github.com/pablor21/goms/pkg/errors"
	"github.com/tidwall/gjson"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

var ErrVideoNotLoaded = errors.NewAppError("Video not loaded", 400)

type VideoMetadata struct {
	Width    int64   `json:"width"`
	Height   int64   `json:"height"`
	Format   string  `json:"format"`
	Codec    string  `json:"codec"`
	Bitrate  int64   `json:"bitrate"`
	Duration float64 `json:"duration"`
}

type PosterParams struct {
	Width  int
	Height int
	Format string
	Time   float64
	Frames int
	Rate   int
	PixFmt string
}

type Video interface {
	Close() error
	GetMetadata() (VideoMetadata, error)
	MakePoster(params PosterParams) (io.ReadWriter, error)
}

type VideoLoader interface {
	Video
	LoadFromReader(reader io.Reader) error
}

type ffmpegVideo struct {
	reader io.Reader
	video  *ffmpeg_go.Stream
}

func NewFFmpegVideo() VideoLoader {
	return &ffmpegVideo{}
}

func (v *ffmpegVideo) Close() error {
	if closer, ok := v.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (v *ffmpegVideo) LoadFromReader(reader io.Reader) error {
	v.reader = reader
	v.video = ffmpeg_go.Input("pipe:0").WithInput(reader)
	return nil
}

func (v *ffmpegVideo) GetMetadata() (res VideoMetadata, err error) {
	if v.reader == nil {
		err = ErrVideoNotLoaded
		return
	}
	// get the video info
	str, err := ffmpeg_go.ProbeReader(v.reader)
	if err != nil {
		return
	}

	// reset the reader
	v.reset()

	res.Bitrate = gjson.Get(str, "streams.0.bit_rate").Int()
	res.Codec = gjson.Get(str, "streams.0.codec_name").String()
	res.Duration = gjson.Get(str, "format.duration").Float()
	res.Format = gjson.Get(str, "format.format_name").String()
	res.Height = gjson.Get(str, "streams.0.height").Int()
	res.Width = gjson.Get(str, "streams.0.width").Int()

	return
}

func (v *ffmpegVideo) MakePoster(params PosterParams) (res io.ReadWriter, err error) {
	res = bytes.NewBuffer(nil)

	if v.reader == nil {
		err = ErrVideoNotLoaded
		return
	}

	if rc, ok := v.reader.(io.ReadSeeker); ok {
		rc.Seek(0, 0)

	}

	err = ffmpeg_go.Input("pipe:0", ffmpeg_go.KwArgs{
		// "ss": params.Time,
	}).
		WithInput(v.reader).
		Output("pipe:1", ffmpeg_go.KwArgs{
			// "vframes": params.Frames,
			"t":       params.Frames,
			"r":       params.Rate,
			"f":       params.Format,
			"pix_fmt": params.PixFmt,
			// "vcodec":  "mgif",
		}).
		// Output("pipe:1", ffmpeg_go.KwArgs{"pix_fmt": "rgb24", "t": "3", "r": "3", "f": "gif"}).
		WithOutput(res).
		Run()
	v.reset()
	return
}

func (v *ffmpegVideo) reset() error {
	if seeker, ok := v.reader.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}
	return nil
}
