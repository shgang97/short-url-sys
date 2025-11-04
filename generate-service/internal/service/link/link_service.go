package link

import (
	"context"
	"fmt"
	"generate-service/internal/model"
	"generate-service/internal/pkg/errors"
	linkRepo "generate-service/internal/repository/link"
	"generate-service/internal/service/idgen"
	"log"
	"time"
)

type linkService struct {
	linkRepo      linkRepo.Repository
	idGenerator   idgen.Generator
	urlValidator  *URLValidator
	codeGenerator *ShortCodeGenerator
	baseURL       string
}

// CreateShortURL 创建短链接
func (s *linkService) CreateShortURL(ctx context.Context, req *model.CreateShortRequest) (*model.Link, error) {
	longURL := req.LongURL
	// 验证URL
	if err := s.ValidateURL(longURL); err != nil {
		return nil, err
	}

	// 标准化URL
	normalizeURL, err := s.NormalizeURL(longURL)
	if err != nil {
		return nil, err
	}

	var shortCode string

	// 处理自定义短码
	if req.CustomCode != nil {
		shortCode = *req.CustomCode

		// 验证自定义短码格式
		if err := s.codeGenerator.ValidateCustomCode(shortCode); err != nil {
			return nil, err
		}

		// 检查短码是否已存在
		exists, err := s.linkRepo.Exists(ctx, shortCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrShortCodeExists
		}
	} else {
		// 生成唯一短码
		var id uint64
		var err error

		// 重试机制，防止ID冲突
		for i := 0; i < 3; i++ {
			id, err = s.idGenerator.NextId()
			if err != nil {
				return nil, err
			}

			shortCode = s.codeGenerator.GenerateFromID(id)

			// 检查短码是否已存在
			exists, err := s.linkRepo.Exists(ctx, shortCode)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}

			// 如果冲突，使用随机短码
			// 不再检查冲突，发生概率极低，即使发生，数据库唯一约束最终保证数据一致性
			if i == 2 {
				shortCode, err = s.codeGenerator.GenerateRandomCode(8)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// 创建链接记录
	createTime := time.Now()
	user := s.getUser(req.CreatedBy)
	id, _ := s.idGenerator.NextId()
	link := &model.Link{
		ID:          id, // 主键ID由雪花算法生成
		ShortCode:   shortCode,
		LongURL:     normalizeURL,
		ExpiresAt:   req.ExpiresAt,
		CreatedBy:   user,
		CreatedAt:   createTime,
		UpdatedBy:   user,
		UpdatedAt:   createTime,
		Status:      model.LinkStatusActive,
		DeleteFlag:  "N",
		Description: s.getDescription(req.Description),
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, err
	}
	// 异步预热缓存
	go func() {
		// TODO 发送缓存预热消息
		log.Printf("发送缓存预热消息")
	}()

	return link, nil
}

// GetLongURL 获取长链接（用于重定向）
func (s *linkService) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	// 直接从数据库获取
	link, err := s.linkRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	// 检查链接状态
	if !link.IsActive() {
		if link.Status == model.LinkStatusExpired {
			return "", errors.ErrLinkExpired
		}
		return "", errors.ErrLinkDisabled
	}

	// 异步更新缓存
	go func() {
		// TODO 发送更新缓存消息
		log.Printf("发送更新缓存消息")
	}()

	return link.LongURL, nil
}

// GetLinkInfo 获取链接信息
func (s *linkService) GetLinkInfo(ctx context.Context, shortCode string) (*model.LinkInfoResponse, error) {
	link, err := s.linkRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	// TODO 远程调用 统计服务获取最后访问时间
	getLastAccess := time.Now()
	lastAccess := &getLastAccess

	linkInfo := &model.LinkInfoResponse{
		ShortCode:    link.ShortCode,
		LongURL:      link.LongURL,
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
		ExpiresAt:    link.ExpiresAt,
		ClickCount:   link.ClickCount,
		LastAccessed: lastAccess,
		Status:       string(link.Status),
		Description:  link.Description,
	}
	return linkInfo, nil
}

// UpdateLink 更新链接信息
func (s *linkService) UpdateLink(ctx context.Context, shortCode string, req *model.UpdateLinkRequest) (*model.LinkInfoResponse, error) {
	link, err := s.linkRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	// 更新字段
	if req.LongURL != nil {
		if err := s.ValidateURL(*req.LongURL); err != nil {
			return nil, err
		}
		normalizeURL, err := s.NormalizeURL(*req.LongURL)
		if err != nil {
			return nil, err
		}
		link.LongURL = normalizeURL
	}
	if req.ExpiresAt != nil {
		link.ExpiresAt = req.ExpiresAt
	}
	if req.Status != nil {
		link.Status = model.LinkStatus(*req.Status)
	}
	if req.Description != nil {
		link.Description = *req.Description
	}
	if err := s.linkRepo.Update(ctx, link); err != nil {
		return nil, err
	}

	// 更新缓存
	go func() {
		if link.Status == model.LinkStatusActive {
			// TODO 发送更新缓存消息
			log.Printf("发送更新缓存消息")
		} else {
			// TODO 发送删除缓存消息
			log.Printf("发送更新缓存消息")
		}
	}()
	// TODO 远程调用 统计服务获取最后访问时间
	getLastAccess := time.Now()
	lastAccess := &getLastAccess

	linkInfo := &model.LinkInfoResponse{
		ShortCode:    link.ShortCode,
		LongURL:      link.LongURL,
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
		ExpiresAt:    link.ExpiresAt,
		ClickCount:   link.ClickCount,
		LastAccessed: lastAccess,
		Status:       string(link.Status),
		Description:  link.Description,
	}
	return linkInfo, nil
}

// DeleteLink 删除链接
func (s *linkService) DeleteLink(ctx context.Context, shortCode string) error {
	// 检查链接是否存在
	link, err := s.linkRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return err
	}
	// TODO 从上下文获取当前用户信息
	link.UpdatedBy = s.getUser(nil)
	if err := s.linkRepo.Delete(ctx, link); err != nil {
		return err
	}
	// 删除缓存
	go func() {
		// TODO 发送删除缓存消息
		log.Printf("发送更新缓存消息")
	}()
	return nil
}

// ListLinks 列表查询链接
func (s *linkService) ListLinks(ctx context.Context, req *model.ListLinksRequest) (*model.ListLinksResponse, error) {
	filter := linkRepo.ListFilter{
		CreatedBy: s.getUser(req.CreatedBy),
		Search:    "",
	}
	if req.Status != nil {
		filter.Status = *req.Status
	}
	links, total, err := s.linkRepo.List(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	// 转换为响应模型
	linkInfos := make([]model.LinkInfoResponse, len(links))
	for i, link := range links {
		// TODO 这里在for循环中查询数据库了，需要优化
		// TODO 远程调用 统计服务获取最后访问时间
		getLastAccess := time.Now()
		lastAccess := &getLastAccess

		linkInfos[i] = model.LinkInfoResponse{
			ShortCode:    link.ShortCode,
			LongURL:      link.LongURL,
			CreatedAt:    link.CreatedAt,
			UpdatedAt:    link.UpdatedAt,
			ExpiresAt:    link.ExpiresAt,
			ClickCount:   link.ClickCount,
			LastAccessed: lastAccess,
			Status:       string(link.Status),
			Description:  link.Description,
		}
	}

	pages := (total + int64(req.PageSize) - 1) / int64(req.PageSize)

	return &model.ListLinksResponse{
		Links:    linkInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    int(pages),
	}, nil
}

// BatchCreate 批量创建链接
func (s *linkService) BatchCreate(ctx context.Context, req *model.BatchCreateRequest) (*model.BatchCreateResponse, error) {
	results := make([]model.BatchResult, 0, len(req.URLs))
	failed := make([]model.BatchFailed, 0, len(req.URLs))

	for _, item := range req.URLs {
		createReq := &model.CreateShortRequest{
			LongURL:    item.LongURL,
			CustomCode: item.CustomCode,
		}

		link, err := s.CreateShortURL(ctx, createReq)
		if err != nil {
			failed = append(failed, model.BatchFailed{
				LongURL: item.LongURL,
				Error:   err.Error(),
			})
			continue
		}
		results = append(results, model.BatchResult{
			LongURL:   link.LongURL,
			ShortURL:  s.buildShortURL(link.ShortCode),
			ShortCode: link.ShortCode,
		})
	}

	return &model.BatchCreateResponse{
		Results: results,
		Failed:  failed,
	}, nil
}

// ValidateURL 验证URL
func (s *linkService) ValidateURL(url string) error {
	return s.urlValidator.Validate(url)
}

// NormalizeURL 标准化URL
func (s *linkService) NormalizeURL(url string) (string, error) {
	return s.urlValidator.NormalizeURL(url)
}

type Config struct {
	BaseURL string
}

// NewService 创建短链服务实例
func NewService(
	linkRepo linkRepo.Repository,
	idGenerator idgen.Generator,
	cfg Config,
) Service {
	return &linkService{
		linkRepo:      linkRepo,
		idGenerator:   idGenerator,
		urlValidator:  NewURLValidator(),
		codeGenerator: NewShortCodeGenerator(),
		baseURL:       cfg.BaseURL,
	}
}

// 获取创建者信息
func (s *linkService) getUser(createdBy *string) string {
	if createdBy == nil {
		return "anonymous"
	}
	return *createdBy
}

// 获取描述信息
func (s *linkService) getDescription(description *string) string {
	if description == nil {
		return ""
	}
	return *description
}

// 构建完整的短链
func (s *linkService) buildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, shortCode)
}
