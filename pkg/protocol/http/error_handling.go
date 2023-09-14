package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcErrorToHttpErrorMap = map[codes.Code]int{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.NotFound:           http.StatusNotFound,
	codes.Aborted:            http.StatusConflict,
	codes.AlreadyExists:      http.StatusConflict,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.Canceled:           499,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unknown:            http.StatusInternalServerError,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
}

func HandleGRPCError(grpcError error) error {
	if grpcError == nil {
		return nil
	}
	if grpcStatus, ok := status.FromError(grpcError); ok {
		detail := extractGrpcErrorDetail(grpcStatus)
		if httpStatus, ok := grpcErrorToHttpErrorMap[grpcStatus.Code()]; ok {
			if detail != "" {
				return echo.NewHTTPError(httpStatus, detail)
			} else {
				return echo.NewHTTPError(httpStatus)
			}
		}
	}
	return echo.ErrInternalServerError
}

func extractGrpcErrorDetail(st *status.Status) string {
	var details bytes.Buffer
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			details.WriteString("Request reject\n")
			for _, violation := range t.GetFieldViolations() {
				_, _ = fmt.Fprintf(&details, "The %q field was wrong:\n", violation.GetField())
				_, _ = fmt.Fprintf(&details, "\t%s\n", violation.GetDescription())
			}
		case *errdetails.ErrorInfo:
			_, _ = fmt.Fprintf(&details, "Domain: %s\n", t.Domain)
			_, _ = fmt.Fprintf(&details, "Metadata: %v\n", t.Metadata)
			_, _ = fmt.Fprintf(&details, "Reason: %s\n", t.Reason)
		case *errdetails.ResourceInfo:
			_, _ = fmt.Fprintf(&details, "Description: %s\n", t.Description)
			_, _ = fmt.Fprintf(&details, "Owner: %s\n", t.Owner)
			_, _ = fmt.Fprintf(&details, "ResourceName: %s\n", t.ResourceName)
			_, _ = fmt.Fprintf(&details, "ResourceType: %s\n", t.ResourceType)
		case *errdetails.DebugInfo:
			_, _ = fmt.Fprintf(&details, "Detail: %s\n", t.Detail)
			for _, stackEntry := range t.StackEntries {
				details.WriteString("\t" + stackEntry + "\n")
			}
		case *errdetails.PreconditionFailure:
			details.WriteString("Precondition Failure\n")
			for _, violation := range t.Violations {
				_, _ = fmt.Fprintf(&details, "- %s (%s): %s\n",
					violation.Subject, violation.Type, violation.Description,
				)
			}
		case *errdetails.QuotaFailure:
			for _, violation := range t.Violations {
				_, _ = fmt.Fprintf(&details, "- %s: %s\n",
					violation.Subject, violation.Description,
				)
			}
		}
	}
	return details.String()
}
