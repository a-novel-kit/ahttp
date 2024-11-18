package ahttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var rpcToHTTPCodes = map[codes.Code]int{
	codes.OK: http.StatusOK,
	// Found it the most appropriate, go does not register this code:
	// https://www.webfx.com/web-development/glossary/http-status-codes/what-is-a-499-status-code/
	codes.Canceled:           499,
	codes.InvalidArgument:    http.StatusUnprocessableEntity,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

func GRPCToHTTPCode(code codes.Code) int {
	if c, ok := rpcToHTTPCodes[code]; ok {
		return c
	}

	return http.StatusInternalServerError
}

// HandleGRPCError handles errors returned by a GRPC service. It returns a boolean indicating
// whether the context was terminated.
func HandleGRPCError(ctx *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	grpcCode, ok := status.FromError(err)
	if !ok {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return true
	}

	_ = ctx.AbortWithError(GRPCToHTTPCode(grpcCode.Code()), grpcCode.Err())
	return true
}
