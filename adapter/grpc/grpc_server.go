package grpcAdapter

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/renatoviolin/shortener/adapter/grpc/pb"
	"github.com/renatoviolin/shortener/application/entity"
	"github.com/renatoviolin/shortener/application/shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	UseCase shortener.UseCaseShortener
}

func NewGrpcHandler(useCase shortener.UseCaseShortener) *GrpcHandler {
	return &GrpcHandler{
		UseCase: useCase,
	}
}

func (h *GrpcHandler) Run(address string) {
	var err error
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShortenerServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	fmt.Printf("GRPC listening on %s\n", address)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC over %s: %v", address, err)
	}
}

func (h *GrpcHandler) GetCode(ctx context.Context, req *pb.Url) (*pb.Code, error) {
	url := req.GetUrl()
	redirect, err := h.UseCase.UrlToCode(url)
	if err != nil {
		if err == entity.ErrRedirectNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	code := &pb.Code{
		Code: redirect.Code,
	}
	return code, nil
}

func (h *GrpcHandler) GetUrl(ctx context.Context, req *pb.Code) (*pb.Url, error) {
	code := req.GetCode()
	redirect, err := h.UseCase.CodeToUrl(code)
	if err != nil {
		if err == entity.ErrRedirectInvalid {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	url := &pb.Url{
		Url: redirect.URL,
	}
	return url, nil
}
