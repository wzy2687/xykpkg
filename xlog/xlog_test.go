package xlog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSlog(t *testing.T) {
	//opts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	customTimeFormat := func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
					//parts := strings.Split(fullPath, ":")
					//if len(parts) > 1 {
					//	fileName = fileName + ":" + parts[len(parts)-1]
					//}
					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(fmt.Sprintf("%s:%d", fileName, source.Line)),
					}
				}
			}
			return a
		},
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	//slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
	slog.Info("info hello ")
	slog.Debug("debug hello ")
	logger.Debug("logger come")
	logger = logger.With("id", "tester")
	logger.Debug("tester's logger")
	logger = logger.WithGroup("grp")
	logger.Debug("grp's logger")
}
