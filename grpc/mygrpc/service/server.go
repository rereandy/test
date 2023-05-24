package service

import (
	"context"
	"fmt"
	pb "mygrpc/proto"
)

type TestServiceImpl struct {
	*pb.UnimplementedTestServiceServer
}

func (t *TestServiceImpl) Query(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	fmt.Println(req.Name)
	resp := &pb.Reply{
		Msg:  "success",
		Code: 0,
	}
	return resp, nil
}
