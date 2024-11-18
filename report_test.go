package ahttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/a-novel-kit/quicklog"
	quicklogmocks "github.com/a-novel-kit/quicklog/mocks"

	"github.com/a-novel-kit/ahttp"
)

func TestReport(t *testing.T) {
	testCases := []struct {
		name string

		status int

		expectLevel quicklog.Level
	}{
		{
			name: "Success",

			status: http.StatusOK,

			expectLevel: quicklog.LevelInfo,
		},
		{
			name: "ClientError",

			status: http.StatusBadRequest,

			expectLevel: quicklog.LevelWarning,
		},
		{
			name: "ServerError",

			status: http.StatusInternalServerError,

			expectLevel: quicklog.LevelError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
			ctx.Writer.WriteHeader(testCase.status)

			logger := quicklogmocks.NewMockLogger(t)
			logger.On("Log", testCase.expectLevel, mock.Anything).Once()

			middleware := ahttp.ReportMiddleware(logger, "hello-world")
			middleware(ctx)

			logger.AssertExpectations(t)
		})
	}
}
