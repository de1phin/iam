package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

var ErrNotFound = errors.New("not found")

type CoreCache interface {
	Get(key string) ([]byte, bool)
	Create(key string, data []byte) bool
	Update(key string, data []byte) bool
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

type TokenKey struct {
	Ssh string
}

type TokenValue struct {
	Token string
}

func (c *MemCache) GetToken(ctx context.Context, ssh string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/GetToken")
	defer span.Finish()

	key, err := marshalStruct(&TokenKey{Ssh: ssh})
	if err != nil {
		return "", err
	}

	value, isExist := c.conn.Get(string(key))
	if !isExist {
		return "", ErrNotFound
	}

	var tokenValue TokenValue
	if err = json.Unmarshal(value, &tokenValue); err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return tokenValue.Token, nil
}

func (c *MemCache) SetToken(ctx context.Context, ssh string, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/SetToken")
	defer span.Finish()

	key, value, err := marshalKeyValue(&TokenKey{Ssh: ssh}, &TokenValue{Token: token})
	if err != nil {
		return err
	}

	isCreated := c.conn.Create(string(key), value)
	if isCreated {
		c.conn.Update(string(key), value)
	}

	return nil
}

func (c *MemCache) DeleteToken(ctx context.Context, ssh string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/DeleteToken")
	defer span.Finish()

	key, err := marshalStruct(&TokenKey{Ssh: ssh})
	if err != nil {
		return err
	}

	c.conn.Delete(string(key))

	return nil
}

type SshKey struct {
	Token string
}

type SshValue struct {
	Ssh string
}

func (c *MemCache) GetSsh(ctx context.Context, token string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/GetSsh")
	defer span.Finish()

	key, err := marshalStruct(&SshKey{Token: token})
	if err != nil {
		return "", err
	}

	value, exist := c.conn.Get(string(key))
	if !exist {
		return "", nil
	}

	var existValue SshValue
	if err = json.Unmarshal(value, &existValue); err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return existValue.Ssh, nil
}

func (c *MemCache) SetSsh(ctx context.Context, token string, ssh string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/SetSsh")
	defer span.Finish()

	key, value, err := marshalKeyValue(&SshKey{Token: token}, &SshValue{Ssh: ssh})
	if err != nil {
		return err
	}

	isCreated := c.conn.Create(string(key), value)
	if isCreated {
		c.conn.Update(string(key), value)
	}

	return nil
}

func (c *MemCache) DeleteSsh(ctx context.Context, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache/DeleteSsh")
	defer span.Finish()

	key, err := marshalStruct(&SshKey{Token: token})
	if err != nil {
		return err
	}

	c.conn.Delete(string(key))

	return nil
}

func marshalKeyValue(key, value any) ([]byte, []byte, error) {
	keyBuff, err := marshalStruct(key)
	if err != nil {
		return nil, nil, err
	}

	valueBuff, err := marshalStruct(value)
	if err != nil {
		return nil, nil, err
	}

	return keyBuff, valueBuff, nil
}

func marshalStruct(str any) ([]byte, error) {
	buff, err := json.Marshal(str)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return buff, nil
}
