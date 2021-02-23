package repositories

import (
	"context"
	"errors"
	"hello-micro/app"
	"hello-micro/form/request"
	"hello-micro/models"
)

type HelloRepositories interface {
	Add(ctx context.Context, req *request.HelloAddReq) error
}

type HelloRepository struct {
}

func (s *HelloRepository) Add(ctx context.Context, req *request.HelloAddReq) error {
	helloModel := models.Hello{
		Name:    req.Name,
		Content: req.Content,
	}
	mysql := app.MysqlDb.Debug()
	if err := mysql.Create(&helloModel).Error; err != nil {
		//日志记录错误
		app.ZapLog.Error("hello-add-error", err.Error())
		return errors.New("create error")
	}
	return nil
}
