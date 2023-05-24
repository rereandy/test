package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "mygrpc/proto"
	"time"
)

// unaryInterceptor 一个简单的 unary interceptor 示例。
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// pre-processing
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing
	end := time.Now()
	log.Printf("RPC: %s, req:%v start time: %s, end time: %s, err: %v", method, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return err
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8002", grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTestServiceClient(conn)
	resp, err := client.Query(context.Background(), &pb.Request{Name: "hello world!"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(resp.Msg)
}
