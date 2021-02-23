package handlers

import (
	"context"
	"hello-micro/app"
	"hello-micro/form/request"
	helloProto "hello-micro/proto/hello"
	"hello-micro/repositories"
)

type HelloHandler struct {
	HelloRepository repositories.HelloRepositories
}

func (s *HelloHandler) Add(ctx context.Context, req *helloProto.AddData, res *helloProto.Response) error {
	data := request.MarshalHelloAdd(req)
	errValidateBool, errValidateMsg := app.GetError(data)
	if errValidateBool == false {
		res.Code = 401
		res.Msg = errValidateMsg
		return nil
	}
	err := s.HelloRepository.Add(ctx, data)
	res.Code = 200
	res.Msg = err.Error()
	return nil
}
