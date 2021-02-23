package main

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	gsrv "github.com/micro/go-micro/v2/server/grpc"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"hello-micro/app"
	"hello-micro/handlers"
	"hello-micro/lib/jaeger"
	helloProto "hello-micro/proto/hello"
	"hello-micro/repositories"
	"log"
	"os"
)

func init() {
	//init config
	app.LoadConfig("./config")
	//init mysql
	app.InitMysql()
	//defer app.MysqlDb.Close()
	//init redis
	app.InitRedis()
	//defer app.RedisDB.Close()

	app.InitLogger()

	if os.Getenv("ENV") == "prod" {
		app.ZapLog.Info("server running", "服务启动 env prod")
	} else {
		app.ZapLog.Info("server running", "服务启动 env test")
	}
}

func main() {
	//链路追踪
	t, closer, err := jaeger.NewJaegerTracer(app.Config.Server.Name, app.Config.Jaeger.Addr)
	if err != nil {
		log.Fatalf("opentracing tracer create error:%v", err)
	}
	defer closer.Close()

	server.DefaultServer = gsrv.NewServer(server.Wait(nil))
	//set up server
	app.MicroService = micro.NewService(
		micro.Name(app.Config.Server.Name),
		micro.Version(app.Config.Server.Version),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(t)), //绑定链路追踪
	)
	app.MicroService.Init()

	repo := &repositories.HelloRepository{}
	h := &handlers.HelloHandler{
		HelloRepository: repo,
	}
	if err := helloProto.RegisterHelloServiceHandler(app.MicroService.Server(), h); err != nil {
		log.Panic(err)
	}
	if err := app.MicroService.Run(); err != nil {
		fmt.Println(err)
	}
}
