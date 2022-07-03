package server

import (
	"fmt"
	"net/http"

	errutils "github.com/mjudeikis/ecowitt-easyweather/pkg/utils/error"
	"github.com/mjudeikis/ecowitt-easyweather/pkg/utils/ratelimiter"
	"go.uber.org/zap"
)

// gRateLimit generic rate limit of unauthenticated endpoints
func (s *Server) gRateLimit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 	// https://cloud.google.com/load-balancing/docs/https
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" { // bad day for you, we don't know who you are so you share rate limit with the rest of the world
			ip = "generic"
		}
		var bailed bool
		if w, bailed = enrichRateLimitOrBail(s.log, w, s.rateLimiters.generic, ip); bailed {
			return
		}
		f(w, r)
	}
}

// enrichRateLimitOrBail returns true if bailed out. False if we are within rate
// limits and enriches responseWriter with data for the rate limiter
func enrichRateLimitOrBail(log *zap.Logger, w http.ResponseWriter, r ratelimiter.Interface, id string) (http.ResponseWriter, bool) {
	wait := r.Wait(id).Seconds()
	available := r.Available(id)
	w.Header().Add("RateLimit-Limit", fmt.Sprintf("%d", r.Cap()))
	w.Header().Add("Ratelimit-Remaining", fmt.Sprintf("%d", available))
	w.Header().Add("Ratelimit-Reset", fmt.Sprintf("%f", wait))

	// if we are under the limit, enrich and return
	if r.Limit(id) {
		log.Debug("rate limit reached", zap.Float64("wait", wait), zap.Int64("available", available), zap.String("identifier", id), zap.String("bucket", r.Name()))
		errutils.WriteCloudError(w, errutils.NewCloudError(http.StatusTooManyRequests, errutils.CloudErrorRateLimit, "rate limit reached, wait %fs", wait))
		return w, true
	}
	return w, false

}
