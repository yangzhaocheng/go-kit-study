package endpoint

import (
	"context"
	"errors"
	"micro-service/rest/domain"
	"micro-service/service"
	"github.com/go-kit/kit/endpoint"
	"strings"
)
var(
	ErrInvalidRequestType =errors.New("request  arith type error")
)
// MakeArithmeticEndpoint make endpoint
func NewArithmeticEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(domain.ArithmeticRequest)

		var (
			res, a, b int
			calError  error
		)

		a = req.A
		b = req.B

		if strings.EqualFold(req.RequestType, "Add") {
			res = svc.Add(a, b)
		} else if strings.EqualFold(req.RequestType, "Substract") {
			res = svc.Subtract(a, b)
		} else if strings.EqualFold(req.RequestType, "Multiply") {
			res = svc.Multiply(a, b)
		} else if strings.EqualFold(req.RequestType, "Divide") {
			res, calError = svc.Divide(a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return domain.ArithmeticResponse{Result: res, Error: calError}, nil
	}
}
