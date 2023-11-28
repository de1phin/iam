package cache

import (
	"context"
	"errors"

	"github.com/de1phin/iam/token/internal/model"
)

var ErrNotFound = errors.New("not found")

type CoreCache interface {
	Get(key string) ([]byte, bool)
	Create(key string, data []byte) bool
	Delete(key string) bool
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

func (c *MemCache) GetToken(ctx context.Context, ssh string) (string, error) {
	return "", nil
}

func (c *MemCache) SetToken(ctx context.Context, ssh string, token string) error {
	return nil
}

func (c *MemCache) GetExist(ctx context.Context, token string) (bool, error) {
	return false, nil
}

func (c *MemCache) SetExist(ctx context.Context, token string, isExist bool) error {
	return nil
}
