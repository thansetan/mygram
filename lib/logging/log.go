package logging

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func New(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("02-Jan-2006 15:04:05 -0700"))
			}
			if a.Key == slog.SourceKey {
				sourceFile := *(a.Value.Any().(*slog.Source))
				split := strings.Split(sourceFile.File, "final-project")
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", split[len(split)-1], sourceFile.Line))
			}
			return a
		},
	}))

	return logger
}
