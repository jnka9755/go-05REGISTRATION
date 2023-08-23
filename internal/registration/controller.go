package registration

import (
	"context"
	"errors"
	"fmt"

	"github.com/jnka9755/go-05META/meta"
	"github.com/jnka9755/go-05RESPONSE/response"

	sdkCourse "github.com/jnka9755/go-05SDKCOURSE/course"
	sdkUser "github.com/jnka9755/go-05SDKUSER/user"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Update Controller
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	GetAllReq struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

	UpdateReq struct {
		ID     string
		Status *string `json:"status"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(b Business, config Config) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(b),
		GetAll: makeGetAllEndpoint(b, config),
		Update: makeUpdateEndpoint(b),
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
			if errors.As(err, &sdkUser.ErrNotFound{}) ||
				errors.As(err, &sdkCourse.ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success create resgistration", responseRegister, nil), nil
	}
}

func makeGetAllEndpoint(b Business, config Config) Controller {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}

		count, err := b.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		registers, err := b.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success get registers", registers, meta), nil
	}
}

func makeUpdateEndpoint(b Business) Controller {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		if err := b.Update(ctx, &req); err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			if errors.As(err, &ErrInvalidStatus{}) {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK(fmt.Sprintf("Success update register with ID -> '%s'", req.ID), nil, nil), nil
	}
}
