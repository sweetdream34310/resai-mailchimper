package constants

import "net/http"

var (
	ErrorInvalidRequest = NewError(http.StatusBadRequest, "invalid request")
	ErrorDatabase       = NewError(http.StatusInternalServerError, "error database")
	ErrorDataNotFound   = NewError(http.StatusNotFound, "data not found")
	ErrorGeneral        = NewError(http.StatusNotImplemented, "general service error")
	ErrorNotAuthorized  = NewError(http.StatusUnauthorized, "user not authorized")
	ErrorEmailNotMatch  = NewError(http.StatusBadRequest, "email not match")

	ErrorDoublePushNotification = NewError(http.StatusBadRequest, "double push notification")
	ErrorClientNotSupported     = NewError(http.StatusUnauthorized, "client not yet configured or present")
)

func NewError(errorCode int, message string) error {
	return &ApplicationError{
		Code:    errorCode,
		Message: message,
	}
}

type ApplicationError struct {
	Code    int
	Message string
}

func (e *ApplicationError) Error() string {
	return e.Message
}
