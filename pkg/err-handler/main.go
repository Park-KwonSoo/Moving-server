package errhandler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorRslt struct {
	RsltCd  string
	RsltMsg string
}

func PanicErr(err error) {
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
}

func AuthorizedErr(req ...interface{}) (ErrorRslt, error) {
	return ErrorRslt{
		RsltCd:  "01",
		RsltMsg: "Authorized",
	}, status.Errorf(codes.Unauthenticated, fmt.Sprintf("%v", req))
}

func NotFoundErr(req ...interface{}) (ErrorRslt, error) {
	return ErrorRslt{
		RsltCd:  "04",
		RsltMsg: "Not Found",
	}, status.Errorf(codes.NotFound, fmt.Sprintf("%v", req))
}

func ConflictErr(req ...interface{}) (ErrorRslt, error) {
	return ErrorRslt{
		RsltCd:  "03",
		RsltMsg: "Conflict",
	}, status.Errorf(codes.AlreadyExists, fmt.Sprintf("%v", req))
}

func ForbiddenErr(req ...interface{}) (ErrorRslt, error) {
	return ErrorRslt{
		RsltCd:  "09",
		RsltMsg: "Forbidden",
	}, status.Errorf(codes.PermissionDenied, fmt.Sprintf("%v", req))
}

func BadRequestErr(req ...interface{}) (ErrorRslt, error) {
	return ErrorRslt{
		RsltCd:  "99",
		RsltMsg: "Bad Request",
	}, status.Errorf(codes.Unknown, fmt.Sprintf("%v", req))
}
