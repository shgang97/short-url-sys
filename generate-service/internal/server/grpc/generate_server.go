package grpc

import (
	"context"
	"generate-service/internal/service/link"
	"log"
	pb "shared/proto/generate"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GenerateServer struct {
	pb.UnimplementedGenerateServiceServer
	linkService link.Service
}

func NewGenerateServer(linkService link.Service) *GenerateServer {
	return &GenerateServer{
		linkService: linkService,
	}
}

func (s *GenerateServer) GetOriginalUrl(ctx context.Context, req *pb.GetOriginalUrlRequest) (*pb.GetOriginalUrlResponse, error) {
	log.Printf("gRPC request received: GetOriginalUrl for %s", req.ShortCode)
	if req.ShortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "short code is required")
	}

	// 调用业务服务
	longUrl, err := s.linkService.GetLongURL(ctx, req.ShortCode)
	if err != nil {
		return nil, status.Error(codes.Canceled, err.Error())
	}
	return &pb.GetOriginalUrlResponse{
		OriginalUrl:  longUrl,
		Exists:       true,
		IsActive:     true,
		ErrorMessage: "",
	}, nil
}
