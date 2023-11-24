package server

import (
	"strings"
	"text/template"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
)

const (
	MATRIX_MUSIC_RECORD_TMPL = `<font color="#FFFF00">[{{if .IsNew}}+{{end}}{{.ID}}]</font> <font color="#009300">{{.Name}} | {{.Band}}</font> <font color="#0000FC">{{.Duration}}</font> {{if .Url}}=&gt; {{.Url}} {{end}}<font color="#FC7F00">{{range $i, $v := .Tags}}{{if $i}} {{end}}#{{$v}}{{end}}</font> <font color="#7F7F7F">[{{.Count}} résultats]</font>`
)

var matrixMusicRecordTmpl = template.Must(template.New("").Parse(MATRIX_MUSIC_RECORD_TMPL))

type MatrixMusicRecordTmplArgs struct {
	IsNew    bool
	ID       int64
	Name     string
	Band     string
	Count    int64
	Tags     []string
	Url      string
	Duration time.Duration
}

type MatrixMusicRecordPrinter struct{}

func (MatrixMusicRecordPrinter) Print(result textbot.Result) string {
	args := MatrixMusicRecordTmplArgs{
		IsNew:    result.IsNew,
		ID:       result.ID,
		Name:     result.Name,
		Band:     result.Band.Name,
		Count:    result.Count,
		Tags:     result.Tags,
		Url:      result.Url,
		Duration: result.Duration,
	}

	var b strings.Builder
	err := matrixMusicRecordTmpl.Execute(&b, args)
	if err != nil {
		panic(err)
	}

	return b.String()
}

type MatrixStatisticsPrinter struct {
	Location *time.Location
}

func (s MatrixStatisticsPrinter) Print(statistics playbot.MusicRecordStatistics) string {
	return ""
}
