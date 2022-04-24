package ports

import (
	"context"

	"github.com/renatoviolin/shortener/adapter/grpc/pb"
	"google.golang.org/grpc"
)

type GRPCPort interface {
	Run() *grpc.Server
	GetCode(context.Context, *pb.Url) (*pb.Code, error)
	GetUrl(context.Context, *pb.Code) (*pb.Url, error)
}
