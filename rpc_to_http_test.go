package ahttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/ahttp"
)

func TestGRPCToHTTPCode(t *testing.T) {
	testCases := []struct {
		name string

		in codes.Code

		expect int
	}{
		{
			name:   "OK",
			in:     codes.OK,
			expect: 200,
		},
		{
			name:   "Canceled",
			in:     codes.Canceled,
			expect: 499,
		},
		{
			name:   "InvalidArgument",
			in:     codes.InvalidArgument,
			expect: 422,
		},
		{
			name:   "DeadlineExceeded",
			in:     codes.DeadlineExceeded,
			expect: 504,
		},
		{
			name:   "NotFound",
			in:     codes.NotFound,
			expect: 404,
		},
		{
			name:   "AlreadyExists",
			in:     codes.AlreadyExists,
			expect: 409,
		},
		{
			name:   "PermissionDenied",
			in:     codes.PermissionDenied,
			expect: 403,
		},
		{
			name:   "ResourceExhausted",
			in:     codes.ResourceExhausted,
			expect: 429,
		},
		{
			name:   "FailedPrecondition",
			in:     codes.FailedPrecondition,
			expect: 400,
		},
		{
			name:   "OutOfRange",
			in:     codes.OutOfRange,
			expect: 400,
		},
		{
			name:   "Unimplemented",
			in:     codes.Unimplemented,
			expect: 501,
		},
		{
			name:   "Unavailable",
			in:     codes.Unavailable,
			expect: 503,
		},
		{
			name:   "Unauthenticated",
			in:     codes.Unauthenticated,
			expect: 401,
		},
		{
			name:   "Default",
			in:     codes.Unknown,
			expect: 500,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if res := ahttp.GRPCToHTTPCode(testCase.in); res != testCase.expect {
				t.Errorf("expected %d, got %d", testCase.expect, res)
			}
		})
	}
}

func TestHandleGRPCError(t *testing.T) {
	testCases := []struct {
		name string

		err error

		expect     bool
		expectCode int
	}{
		{
			name: "NilError",

			expect:     false,
			expectCode: http.StatusOK,
		},
		{
			name: "NilStatus",

			err: testutils.ErrDummy,

			expect:     true,
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "StatusError",

			err: status.Error(codes.NotFound, "foo"),

			expect:     true,
			expectCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)

			ok := ahttp.HandleGRPCError(ctx, testCase.err)
			require.Equal(t, testCase.expect, ok)
			require.Equal(t, testCase.expectCode, w.Code)
		})
	}
}
