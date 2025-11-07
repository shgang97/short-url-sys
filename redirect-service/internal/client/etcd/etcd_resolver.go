package etcdresolver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdResolver struct {
	etcdClient *clientv3.Client
	conn       resolver.ClientConn
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	config     *EtcdConfig
	key        string
}

type EtcdConfig struct {
	Endpoints []string       `mapstructure:"endpoints"` // ETCD集群地址
	Resolver  ResolverConfig `mapstructure:"resolver"`
}

type ResolverConfig struct {
	ServiceName string        `mapstructure:"service_name"` // 服务名
	Ttl         int64         `mapstructure:"ttl"`          // 租约存活时间
	DialTimeout time.Duration `mapstructure:"dial_timeout"` // ETCD连接超时时间
	Prefix      string        `mapstructure:"prefix"`
}

// Build 创建解析器
func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 解析 target.URL.Path，格式：/serviceName
	serviceName := strings.TrimPrefix(target.URL.Path, "/")
	if serviceName == "" {
		return nil, errors.New("empty service name")
	}
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   r.config.Endpoints,
		DialTimeout: r.config.Resolver.DialTimeout * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	resolve := &etcdResolver{
		etcdClient: etcdClient,
		conn:       cc,
		ctx:        ctx,
		cancel:     cancel,
		key:        fmt.Sprintf("/%s/%s", r.config.Resolver.Prefix, serviceName), // /services/generate-service
	}

	// 启动服务发现
	resolve.wg.Add(1)
	go resolve.watcher()

	return resolve, nil
}

func (r *etcdResolver) watcher() {
	defer r.wg.Done()

	// 初次获取所有服务实例
	r.updateState()

	// 监听服务变化
	watchChan := r.etcdClient.Watch(context.Background(), r.key, clientv3.WithPrefix())

	for {
		select {
		case <-r.ctx.Done():
			return
		case resp := <-watchChan:
			if resp.Canceled {
				log.Printf("etcd watcher got cancelled")
				return
			}
			for range resp.Events {
				r.updateState()
			}
		}
	}

}

// Scheme 返回解析器的scheme
func (r *etcdResolver) Scheme() string {
	return "etcd"
}

func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.cancel()
	r.wg.Wait()
	r.etcdClient.Close()
}

// 更新服务地址列表
func (r *etcdResolver) updateState() {
	resp, err := r.etcdClient.Get(r.ctx, r.key, clientv3.WithPrefix())
	if err != nil {
		log.Printf("ETCD get failed: %v", err)
		return
	}

	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		addr := string(kv.Value)
		if addr != "" {
			addrs = append(addrs, resolver.Address{Addr: addr})
			log.Printf("Discovered service instance: %s", addr)
		}
	}

	if len(addrs) == 0 {
		log.Printf("ETCD service instance not found")
		return
	}

	// 更新 gRPC 客户端连接状态
	r.conn.UpdateState(resolver.State{Addresses: addrs})
	log.Printf("Updaate service address: %v", addrs)
}

// Builder 创建解析器 Builder
type Builder struct {
	config *EtcdConfig
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &etcdResolver{config: b.config}
	return r.Build(target, cc, opts)
}

// Scheme 返回scheme
func (b *Builder) Scheme() string {
	return "etcd"
}

// Init 初始化解析器
func Init(cfg *EtcdConfig) {
	resolver.Register(&Builder{config: cfg})
}
