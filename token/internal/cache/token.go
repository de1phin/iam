package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (c *MemCache) GetToken(_ context.Context, ssh string) (string, error) {
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

func (c *MemCache) SetToken(_ context.Context, ssh string, token string) error {
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

type ExistKey struct {
	Token string
}

type ExistValue struct {
	IsExist bool
}

func (c *MemCache) GetExist(_ context.Context, token string) (bool, error) {
	key, err := marshalStruct(&ExistKey{Token: token})
	if err != nil {
		return false, err
	}

	value, exist := c.conn.Get(string(key))
	if !exist {
		return false, nil
	}

	var existValue ExistValue
	if err = json.Unmarshal(value, &existValue); err != nil {
		return false, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return existValue.IsExist, nil
}

func (c *MemCache) SetExist(_ context.Context, token string, isExist bool) error {
	key, value, err := marshalKeyValue(&ExistKey{Token: token}, &ExistValue{IsExist: isExist})
	if err != nil {
		return err
	}

	isCreated := c.conn.Create(string(key), value)
	if isCreated {
		c.conn.Update(string(key), value)
	}

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
