package log

import (
	"github.com/mjudeikis/ecowitt-easyweather/pkg/api"
	"go.uber.org/zap"
)

// EnrichWithCorrelationData sets log fields based on an optional
// correlationData struct
func EnrichWithCorrelationData(rlog *zap.Logger, correlationData *api.CorrelationData) *zap.Logger {
	if correlationData == nil {
		return rlog
	}

	return rlog.With(
		zap.String("request_id", correlationData.RequestID),
	)
}
