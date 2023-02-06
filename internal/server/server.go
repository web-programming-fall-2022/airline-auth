package server

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "github.com/web-programming-fall-2022/airline-auth/pkg/api/v1"
	"google.golang.org/grpc/metadata"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer() *AuthServiceServer {
	return &AuthServiceServer{}
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, _ = metadata.FromIncomingContext(ctx)
	logrus.Info(ctx.Value("Authentication"))
	return &pb.LoginResponse{}, nil
}
