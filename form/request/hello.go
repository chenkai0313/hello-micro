package request

import helloProto "hello-micro/proto/hello"

type HelloAddReq struct {
	Name    string `json:"name" form:"name" validate:"required"`
	Content string `json:"content" form:"content" validate:"required"`
}

func MarshalHelloAdd(addData *helloProto.AddData) *HelloAddReq {
	return &HelloAddReq{
		Name:    addData.Name,
		Content: addData.Content,
	}
}
