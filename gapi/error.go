package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolations(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid request")
	details, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return details.Err()
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
}

func permissionDeniedError(err error) error {
	return status.Errorf(codes.PermissionDenied, "permission denied: %v", err)
}
