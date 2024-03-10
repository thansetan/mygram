package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type wrappedRW struct {
	code int
	http.ResponseWriter
}

func (wRW *wrappedRW) WriteHeader(code int) {
	wRW.code = code
	wRW.ResponseWriter.WriteHeader(code)
}

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("02-Jan-2006 15:04:05 -0700"))
		}
		return a
	},
}))

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wRW := &wrappedRW{
			ResponseWriter: w,
		}
		t0 := time.Now()
		next.ServeHTTP(wRW, r)
		if wRW.code < 400 {
			logger.Info(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", wRW.code, "took", time.Since(t0))
		} else if wRW.code < 500 {
			logger.Warn(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", wRW.code, "took", time.Since(t0))
		} else {
			logger.Error(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", wRW.code, "took", time.Since(t0))
		}
	})
}
