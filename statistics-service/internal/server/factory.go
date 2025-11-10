package server

import (
	"fmt"
	"shared/constants"
	"statistics-service/internal/consumer"
)

// 资源初始化顺序
// 基础工具中间件：日志中间件、ID生成器
// 数据库：MySQL、Redis
// Repository层
// Service层
// Kafka消费者
// 路由及Handler层
// HTTP服务器

// MessageHandler创建工厂
func (s *Server) createMessageHandler(groupId, topic string, cfg *consumer.GroupConfig) (consumer.MessageHandler, error) {
	handlerKey := consumer.GetHandlerKey(groupId, topic)
	switch handlerKey {
	case consumer.GetHandlerKey(constants.StatsGroupDetail, constants.TopicRecordClickEvent):
		return consumer.NewRecordClickHandler(s.clickSvc), nil
	case consumer.GetHandlerKey(constants.StatsGroupTotal, constants.TopicRecordClickEvent):
		return consumer.NewSummaryHandler(handlerKey, s.summarySvc, cfg.BatchSize, cfg.Spec), nil
	default:
		return nil, fmt.Errorf("unknown topic '%s'", topic)
	}
}
