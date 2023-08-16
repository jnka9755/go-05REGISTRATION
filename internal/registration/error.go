package registration

import (
	"errors"
)

var ErrUserIDRequired = errors.New("user_id is required")
var ErrCourseIDRequired = errors.New("course_id is required")
