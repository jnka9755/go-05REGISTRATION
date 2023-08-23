package registration

import (
	"errors"
	"fmt"
)

var ErrUserIDRequired = errors.New("user_id is required")
var ErrCourseIDRequired = errors.New("course_id is required")
var ErrStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	RegisterID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Register with ID -> '%s' doesn't exist", e.RegisterID)
}

type ErrInvalidStatus struct {
	Status string
}

func (e ErrInvalidStatus) Error() string {
	return fmt.Sprintf("Ivalid '%s' status", e.Status)
}
