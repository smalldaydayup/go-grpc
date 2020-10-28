package grpc

import (
	"context"
	v1 "github.comgo-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// 註冊服務
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	// 優雅地關閉
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// 信號是CTRL+C
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
		}
	}()

	// 啟動gRPC服務器
	log.Println("starting gRPC server...")
	// 原作者的啟動gRPC服務器是這樣子的，但我覺得不太好，所以我改為我的方式去啟動
	// return server.Serve(listen)
	if err := server.Serve(listen); err != nil {
		log.Fatal("starting gRPC server failed...")
		return err
	}

	return nil
}
