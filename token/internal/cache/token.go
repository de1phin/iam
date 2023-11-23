package cache

import (
	"context"

	"github.com/de1phin/iam/token/internal/model"
)

type CoreCache interface {
}

type MemCache struct {
	conn CoreCache
}

func NewCache(conn CoreCache) *MemCache {
	return &MemCache{
		conn: conn,
	}
}

func (c *MemCache) GenerateToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (c *MemCache) RefreshToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (c *MemCache) DeleteToken(ctx context.Context) error {
	return nil
}

func (c *MemCache) CheckToken(context.Context) (bool, error) {
	return false, nil
}
