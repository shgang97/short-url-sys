package register

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceRegister struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

type EtcdConfig struct {
	Endpoints []string       `mapstructure:"endpoints"` // ETCD集群地址
	Register  RegisterConfig `mapstructure:"register"`
}

type RegisterConfig struct {
	ServiceName string        `mapstructure:"service_name"` // 服务名
	Addr        string        `mapstructure:"addr"`         // 当前服务实例的gRPC服务地址
	Ttl         int64         `mapstructure:"ttl"`          // 租约存活时间
	DialTimeout time.Duration `mapstructure:"dial_timeout"` // ETCD连接超时时间
}

func NewServiceRegister(cfg *EtcdConfig) (*ServiceRegister, error) {
	// 创建 ETCD 客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.Register.DialTimeout * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}
	reg := &ServiceRegister{
		cli: cli,
		key: fmt.Sprintf("/services/%s/%s", cfg.Register.ServiceName, cfg.Register.Addr),
		val: cfg.Register.Addr,
	}

	// 执行服务注册（设置键值对并绑定租约）
	err = reg.putKeyWithLease(cfg.Register.Ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to put key with lease: %w", err)
	}

	return reg, nil
}

// 设置键值对并绑定租约
func (s *ServiceRegister) putKeyWithLease(ttl int64) error {
	// 创建租约（设置TTL）
	resp, err := s.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	// 注册服务：将服务信息写入ETCD，并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 启动自动续租机制
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	// 保存租约相关信息
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan

	log.Printf("Service register: %s -> %s with leaseID: %d", s.key, s.val, s.leaseID)
	return nil
}

// ListenKeepAlive 监听续租情况
func (s *ServiceRegister) ListenKeepAlive(ttl int64) {
	for {
		select {
		case resp := <-s.keepAliveChan:
			if resp == nil {
				// 续租通道关闭，说明续租失败
				log.Printf("Keep alive channel closed, trying to re-register...")
				if err := s.putKeyWithLease(ttl); err != nil {
					// 尝试重新注册服务
					log.Printf("failed to re-register key with lease: %v", err)
				}
				return
			}
			log.Printf("KeepAlive success: %v", resp)
		}
	}
}

// Close 关闭服务注册
func (s *ServiceRegister) Close() error {
	// 撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Printf("Service lease revoked: %d", s.leaseID)
	return s.cli.Close()
}
