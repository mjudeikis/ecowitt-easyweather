package recover

import (
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"

	promutil "github.com/mjudeikis/ecowitt-easyweather/pkg/utils/prometheus"
)

// Panic recovers a panic
func Panic(log *zap.Logger) {
	if e := recover(); e != nil {
		log.Error(fmt.Sprint("%w", e))
		log.Info(string(debug.Stack()))

		promutil.GOPanicCounter.WithLabelValues(string(debug.Stack())).Inc()
	}
}
