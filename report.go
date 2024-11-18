package ahttp

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/a-novel-kit/quicklog"

	ahttpmessages "github.com/a-novel-kit/ahttp/messages"
)

func ReportMiddleware(logger quicklog.Logger, projectID string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		level := quicklog.LevelInfo
		if ctx.Writer.Status() >= 500 {
			level = quicklog.LevelError
		} else if ctx.Writer.Status() >= 400 {
			level = quicklog.LevelWarning
		}

		logger.Log(level, ahttpmessages.NewReport(
			&ahttpmessages.Metrics{
				Latency:   time.Since(start),
				StartedAt: start,
			},
			projectID,
			ctx,
		))
	}
}
