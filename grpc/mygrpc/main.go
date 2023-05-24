package main

import (
	"context"
	grpc "google.golang.org/grpc"
	"log"
	pb "mygrpc/proto"
	"mygrpc/service"
	"net"
	"time"
)

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	log.Printf("RPC: %s,req:%v start time: %s, end time: %s, err: %v", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}

func main() {
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptor))
	pb.RegisterTestServiceServer(srv, &service.TestServiceImpl{})
	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = srv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
