package xlog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func GetASlogLogger(w io.Writer) *slog.Logger {

	customTimeFormat := func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05.000000")
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(customTimeFormat(a.Value.Time())),
				}
			}
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					fullPath := source.File
					fileName := filepath.Base(fullPath)
					fn := source.Function
					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(fmt.Sprintf("%s:%d fn:%s", fileName, source.Line, fn)),
					}
				}
			}
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	//slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return logger
}

func init() {
	InitSlog()
}

func InitSlog() {
	slog.SetDefault(GetASlogLogger(os.Stdout))
}

// InitSlog2File
// logName: no date
// path refer env:_go_log_path
func InitSlogFileWithRotateDay(logName string) {

	if len(logName) == 0 {
		_, logName = filepath.Split(os.Args[0])
	}

	fPath := os.Getenv("_go_log_path") //"./logs/"
	if len(fPath) == 0 {
		fPath = "./logs/"
	}
	logFullPath := filepath.Join(fPath, logName)

	l := newWriter(logFullPath)
	write2Stdout := os.Getenv("_go_log2stdout") == "1"
	var w io.Writer
	if write2Stdout {
		w = io.MultiWriter(l, os.Stdout)
	} else {
		w = io.MultiWriter(l)
	}
	slog.SetDefault(GetASlogLogger(w))
}

func MustNoErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "err", err.Error())
		panic(err)
	}
}
