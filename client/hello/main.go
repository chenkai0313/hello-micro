package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	helloProto "hello-micro/proto/hello"
)

const serverName = "go.micro.srv.hello"

type HelloService struct {
}

func main() {
	service := micro.NewService(micro.Name(serverName))
	service.Init()
	client := helloProto.NewHelloService(serverName, service.Client())
	r, err := client.Add(context.Background(), &helloProto.AddData{
		Name:    "hello",
		Content: "tom",
	})
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(r.Msg)
}
