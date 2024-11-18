package ahttpmessages_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	testutils "github.com/a-novel-kit/test-utils"

	ahttpmessages "github.com/a-novel-kit/ahttp/messages"
)

func TestReport(t *testing.T) {
	testCases := []struct {
		name string

		metrics   *ahttpmessages.Metrics
		projectID string
		ginC      func() *gin.Context

		expect     string
		expectJSON map[string]interface{}
	}{
		{
			name: "SimpleRequest",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo]\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "INFO",
			},
		},
		{
			name: "5XX code",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusInternalServerError)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "👶🔪🩸 500 [GET /foo]\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        500,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "ERROR",
			},
		},
		{
			name: "4XX code",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusBadRequest)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "⚠ 400 [GET /foo]\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        400,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "WARNING",
			},
		},
		{
			name: "WithMetrics",

			metrics: &ahttpmessages.Metrics{
				StartedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Latency:   time.Second,
			},
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo] (1s)\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
					"latency":       "1s",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "INFO",
				"start":    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "WithQuery",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo?foo=bar", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo]\n" +
				"┌───────────────────────────────────────┬──────────────────────────────────────┐\n" +
				"│  foo                                  │  bar                                 │\n" +
				"└───────────────────────────────────────┴──────────────────────────────────────┘\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{"foo": []string{"bar"}},
				"severity": "INFO",
			},
		},
		{
			name: "WithQuery",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo?foo=bar&bar=baz&foo=qux", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo]\n" +
				"┌───────────────────────────────────────┬──────────────────────────────────────┐\n" +
				"│  bar                                  │  baz                                 │\n" +
				"│  foo                                  │  bar                                 │\n" +
				"│                                       │  qux                                 │\n" +
				"└───────────────────────────────────────┴──────────────────────────────────────┘\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
				},
				"ip": "127.0.0.1",
				"query": url.Values{
					"bar": []string{"baz"},
					"foo": []string{"bar", "qux"},
				},
				"severity": "INFO",
			},
		},
		{
			name: "WithErrors",

			metrics:   nil,
			projectID: "",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusInternalServerError)

				_ = ctx.Error(errors.New("if you pour water on a rock"))
				_ = ctx.Error(errors.New("but if you pour water on a rock"))

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "👶🔪🩸 500 [GET /foo]\n" +
				"  - if you pour water on a rock\n" +
				"  - but if you pour water on a rock\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors": []string{
					"if you pour water on a rock",
					"but if you pour water on a rock",
				},
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        500,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "ERROR",
			},
		},
		{
			name: "WithTrace",

			metrics:   nil,
			projectID: "cd",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("X-Cloud-Trace-Context", "abcdefg/hijklmnop")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo]\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
				},
				"ip":                           "127.0.0.1",
				"query":                        url.Values{},
				"severity":                     "INFO",
				"logging.googleapis.com/trace": "projects/cd/traces/abcdefg",
			},
		},
		{
			name: "WithTrace/NoTrace",

			metrics:   nil,
			projectID: "cd",
			ginC: func() *gin.Context {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)
				ctx.Request.Header.Set("User-Agent", "Netscape")
				ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Writer.WriteHeader(http.StatusOK)

				_ = testutils.AssignPrivateField[gin.Context, string](ctx, "fullPath", "/foo")

				return ctx
			},

			expect: "✅ 200 [GET /foo]\n\n",
			expectJSON: map[string]interface{}{
				"contentType": "application/json",
				"errors":      []string(nil),
				"httpRequest": map[string]interface{}{
					"protocol":      "HTTP/1.1",
					"remoteIp":      "127.0.0.1",
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
				},
				"ip":       "127.0.0.1",
				"query":    url.Values{},
				"severity": "INFO",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			report := ahttpmessages.NewReport(testCase.metrics, testCase.projectID, testCase.ginC())
			require.Equal(t, testCase.expect, report.RenderTerminal())
			require.Equal(t, testCase.expectJSON, report.RenderJSON())
		})
	}
}
