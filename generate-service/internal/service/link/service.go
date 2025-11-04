package link

import (
	"context"
	"generate-service/internal/model"
)

// Service 短链服务接口
type Service interface {
	CreateShortURL(ctx context.Context, req *model.CreateShortRequest) (*model.Link, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	GetLinkInfo(ctx context.Context, shortCode string) (*model.LinkInfoResponse, error)
	UpdateLink(ctx context.Context, shortCode string, req *model.UpdateLinkRequest) (*model.LinkInfoResponse, error)
	DeleteLink(ctx context.Context, shortCode string) error
	ListLinks(ctx context.Context, req *model.ListLinksRequest) (*model.ListLinksResponse, error)
	BatchCreate(ctx context.Context, req *model.BatchCreateRequest) (*model.BatchCreateResponse, error)
	ValidateURL(url string) error
	NormalizeURL(url string) (string, error)
}
