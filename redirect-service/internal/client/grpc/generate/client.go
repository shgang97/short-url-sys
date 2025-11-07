package generate

import (
	"context"
	"fmt"
	"log"
	"redirect-service/internal/client/etcd"
	"redirect-service/internal/config"
	pb "shared/proto/generate"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.GenerateServiceClient
	config *config.GenerateService
}

// NewClient 创建 gRPC 客户端
func NewClient(cfg *config.GenerateService, etcdCfg *etcdresolver.EtcdConfig) (*Client, error) {
	// 初始化并注册ETCD解析器
	etcdresolver.Init(etcdCfg)
	// gRPC连接选项
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(10*1024*1024),
			grpc.MaxCallRecvMsgSize(10*1024*1024),
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 启用轮询负载均衡
	}
	// 建立连接
	//conn, err := grpc.NewClient(cfg.Address, dialOptions...)
	target := fmt.Sprintf("etcd///%s", etcdCfg.Resolver.ServiceName)
	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}
	client := pb.NewGenerateServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
		config: cfg,
	}, nil
}

func (c *Client) GetOriginalURL(ctx context.Context, shortCode string) (*pb.GetOriginalUrlResponse, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.GetOriginalUrl(ctx, &pb.GetOriginalUrlRequest{
		ShortCode: shortCode,
	})
	if err != nil {
		log.Printf("gRPC call failed for %s: %v", shortCode, err)
		return nil, err
	}
	return resp, nil
}
