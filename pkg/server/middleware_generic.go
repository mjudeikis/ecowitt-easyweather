package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"

	"github.com/mjudeikis/weewx-easyweather/pkg/api"
	errutils "github.com/mjudeikis/weewx-easyweather/pkg/utils/error"
	logutil "github.com/mjudeikis/weewx-easyweather/pkg/utils/log"
	promutil "github.com/mjudeikis/weewx-easyweather/pkg/utils/prometheus"
)

type contextKey int

const (
	ContextKeyLog contextKey = iota
	ContextKeyOriginalPath
	ContextKeyBody
	ContextKeyCorrelationData

	ContextKeyUser
	ContextKeyDevice
	ContextKeyProject
	ContextKeyNamespace
	ContextKeyServiceAccount
)

func Panic(log *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					log.Error("panic")
					log.Sugar().Error(e)

					promutil.GOPanicCounter.WithLabelValues(string(debug.Stack())).Inc()

					log.Sugar().Error(string(debug.Stack()))
					errutils.WriteCloudError(w, errutils.NewCloudError(http.StatusInternalServerError, errutils.CloudErrorCodeInternalServerError, ""))
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

type logResponseWriter struct {
	http.ResponseWriter

	statusCode int
	path       string
	bytes      int
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *logResponseWriter) Flush() {
	flucher := w.ResponseWriter.(http.Flusher)
	flucher.Flush()
}

func (w *logResponseWriter) CloseNotify() <-chan bool {
	notify := w.ResponseWriter.(http.CloseNotifier)
	return notify.CloseNotify()
}

type logReadCloser struct {
	io.ReadCloser

	bytes int
}

func (rc *logReadCloser) Read(b []byte) (int, error) {
	n, err := rc.ReadCloser.Read(b)
	rc.bytes += n
	return n, err
}

func Log(log *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			timer := prometheus.NewTimer(promutil.HttpDuration.WithLabelValues(path, r.Method))
			timer.ObserveDuration()

			r.Body = &logReadCloser{ReadCloser: r.Body}
			w = &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK, path: r.URL.Path}

			correlationData := &api.CorrelationData{
				RequestID:   uuid.NewV4().String(),
				RequestTime: t,
			}

			rlog := log
			rlog = logutil.EnrichWithCorrelationData(rlog, correlationData)

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyLog, rlog)
			ctx = context.WithValue(ctx, ContextKeyCorrelationData, correlationData)

			r = r.WithContext(ctx)

			rlog = rlog.With(
				zap.String("request_method", r.Method),
				zap.String("request_path", r.URL.Path),
				zap.String("request_proto", r.Proto),
				zap.String("request_remote_addr", r.RemoteAddr),
				zap.String("request_user_agent", r.UserAgent()),
			)

			defer func() {
				promutil.TotalRequests.WithLabelValues(
					path,
					strconv.Itoa(w.(*logResponseWriter).statusCode),
					r.Method,
					r.UserAgent(),
				).Inc()

				promutil.BytesReceivedCounter.WithLabelValues(path).Add(float64(r.Body.(*logReadCloser).bytes))
				promutil.BytesTransferredCounter.WithLabelValues(path).Add(float64(w.(*logResponseWriter).bytes))

				rlog.With(
					zap.Int("body_read_bytes", r.Body.(*logReadCloser).bytes),
					zap.Int("body_written_bytes", w.(*logResponseWriter).bytes),
					zap.Float64("duration", time.Since(t).Seconds()),
					zap.Int("response_status_code", w.(*logResponseWriter).statusCode),
				).Warn("sent response")
			}()
			h.ServeHTTP(w, r)
		})
	}
}

func Gzip() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.CompressHandler(h)
	}
}

// shouldLog defines if we should log request.
// Current rules:
// 1. If we returning an error (>=399) - check. Else not log
// 2. If 401 (unauthorized) - don't log
// 3. If request originated not from our up tp date agent/cli/ui - don't log
func shouldLog(w http.ResponseWriter) bool {
	statusCode := w.(*logResponseWriter).statusCode

	// TODO: Once agent checks are rollout and we don't see many default agent ("Go-http-client/2.0)
	// metrics - drop all these
	if statusCode >= 399 {
		// we don't log unauth as they are noisy
		if statusCode == 401 {
			return false
		}

	}
	return false
}
