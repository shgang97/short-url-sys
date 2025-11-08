# 高并发分布式短链系统

## 项目简介

基于Go语言开发的高性能、分布式URL缩短服务，采用微服务架构设计，支持海量短链生成、毫秒级重定向和实时数据统计分析。系统具备完整的监控告警和容器化部署能力，可支撑日均千万级请求。

## 核心特性

- 🚀 **高性能设计**：毫秒级重定向响应，支撑万级QPS并发访问
- 🔗 **微服务架构**：服务解耦，独立扩展，高可用保障
- 📊 **实时统计**：完整的点击数据统计与分析
- 🔄 **异步消息驱动**：Kafka消息队列，系统解耦与流量削峰
- 📈 **全方位监控**：Prometheus + Grafana 实时监控
- 🐳 **容器化部署**：Docker容器化，支持弹性扩缩容
- ⚖️ **负载均衡**：gRPC + ETCD 服务发现与负载均衡

## ⏳待开发功能
- [ ] 完整的点击数据统计与分析
- [ ] Prometheus + Grafana 实时监控

## 系统架构

### 架构概述

系统采用微服务架构，包含三个核心服务：

```
┌────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│      短链生成服务    │    │      重定向服务       │    │     点击统计服务      │
│ (Generate Service) │    │ (Redirect Service)  │    │ (Statistics Service)│
└────────────────────┘    └─────────────────────┘    └─────────────────────┘
         │                           │                           │
         ├──────────Kafka───────────▶│                          │
         │                           ├──────────Kafka───────────▶│
         │◀︎─────────gRPC─────────────│                            │
┌────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│     MySQL          │    │      Redis          │    │    Prometheus       │
│   (数据持久化)       │    │   (缓存层)           │    │    (监控指标)        │
└────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

### 核心服务

1. **短链生成服务**
    - 短链的创建、更新、删除
    - 发送Kafka消息通知缓存更新
    - 提供gRPC接口供其他服务调用

2. **重定向服务**
    - 处理短链访问请求
    - Redis缓存优先，缓存未命中时gRPC回溯
    - 发送点击事件到Kafka

3. **点击统计服务**
    - 消费点击事件消息
    - 统计分析与数据聚合
    - 提供统计查询API

## 技术栈

### 后端技术
- **编程语言**: Go
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: MySQL
- **缓存**: Redis
- **消息队列**: Kafka
- **服务通信**: gRPC
- **服务发现**: ETCD

### 运维监控
- **容器化**: Docker + Docker Compose
- **监控**: Prometheus + Grafana
- **指标采集**: 各服务自定义指标

## 快速开始

### 环境要求
- Docker & Docker Compose
- Go 1.25

### 部署步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/shgang97/short-url-sys.git
   cd short-url-sys
   ```

2. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，配置数据库、Redis等连接信息
   ```

3. **启动服务**
   ```bash
   docker-compose up -d
   ```

4. **验证部署**
   ```bash
   # 检查服务状态
   docker-compose ps
   
   # 查看服务日志
   docker-compose logs -f [service-name]
   ```

### 访问地址
- **API服务**: http://localhost:8080
- **监控看板**: http://localhost:3000 (Grafana)
- **指标接口**: http://localhost:9090 (Prometheus)

## API文档

### 生成短链
```http
POST /api/v1/links/short
Content-Type: application/json

{
  "long_url": ""
}
```

### 访问短链
```http
GET /:code
```

### 获取统计信息
```http
GET /api/v1/stats/:code
```

## 性能指标

- **重定向延迟**: P99 < 25ms
- **生成服务TPS**: 
- **缓存命中率**: 
- **系统吞吐量**: 

## 监控指标

系统监控以下关键指标：
- API响应时间与QPS
- Kafka消息消费延迟
- Redis缓存命中率
- MySQL查询性能
- 服务实例健康状态

## 项目结构

```
short-url-sys/
├── generate-service/
│   ├── cmd/
│   ├── configs/ 
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── pkg/
│   │   ├── repository/
│   │   ├── server/
│   │   └── service/
│   ├── scripts/
│   ├── docker-compose.yml 
│── └── go.mod
│── redirect-service/ 
│   ├── cmd/
│   ├── configs/ 
│   ├── data/
│   ├── internal/
│   │   ├── client/
│   │   ├── config/
│   │   ├── consumer/
│   │   ├── handler/
│   │   ├── pkg/
│   │   ├── producer/
│   │   ├── repository/
│   │   ├── server/
│   │   └── service/
│── └── go.mod
│── statistics-service/
│   ├── cmd/
│   ├── configs/ 
│   ├── internal/
│   │   ├── config/
│   │   ├── consumer/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── pkg/
│   │   ├── repository/
│   │   ├── server/
│   │   └── service/
│   ├── scripts/
│── └── go.mod
│── shared/
│   ├── constants/
│   ├── errors/ 
│   ├── message/
│   ├── model/
│   ├── proto/
│── └── go.mod
│── go.work
│── LICENSE
└── README.md
```

## 开发指南

### 本地开发
```bash
# 启动依赖服务
docker-compose up -d mysql redis kafka etcd

# 运行生成服务
cd cmd/generate-service && go run main.go

# 运行重定向服务  
cd cmd/redirect-service && go run main.go

# 运行统计服务
cd cmd/analytics-service && go run main.go
```

### 构建镜像
```bash
# 构建所有服务镜像
docker-compose build

# 构建特定服务
docker-compose build generate-service
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目主页: https://github.com/shgang97/short-url-sys
- 问题反馈: https://github.com/shgang97/short-url-sys/issues

---

**星星这个项目 ⭐ 如果你觉得这个项目对你有帮助！**