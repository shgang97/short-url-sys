package grpc

import (
	"context"
	"generate-service/internal/service/link"
	"log"
	pb "shared/proto/generate"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	lk, err := s.linkService.GetLink(ctx, req.ShortCode)
	if err != nil {
		return nil, status.Error(codes.Canceled, err.Error())
	}
	resp := &pb.GetOriginalUrlResponse{
		OriginalUrl:  lk.LongURL,
		IsActive:     true,
		ErrorMessage: "",
	}
	if lk.ExpiresAt != nil && !lk.ExpiresAt.IsZero() {
		resp.ExpireTime = timestamppb.New(*lk.ExpiresAt)
	}
	return resp, nil
}
