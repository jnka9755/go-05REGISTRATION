package registration

import (
	"context"

	"github.com/jnka9755/go-05RESPONSE/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}
)

func MakeEndpoints(b Business) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(b),
	}
}

func makeCreateEndpoint(b Business) Controller {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.UserID == "" {
			return nil, response.BadRequest(ErrUserIDRequired.Error())
		}

		if req.CourseID == "" {
			return nil, response.BadRequest(ErrCourseIDRequired.Error())
		}

		responseRegister, err := b.Create(ctx, &req)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success create resgistration", responseRegister, nil), nil
	}
}
