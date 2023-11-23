package cache

import (
	account "github.com/de1phin/iam/genproto/services/account/api"
	pkg_cache "github.com/de1phin/iam/pkg/cache"
	"github.com/de1phin/iam/services/account/internal/database"
)

type AccountCache struct {
	*pkg_cache.Cache[string, *account.Account]
}

func NewAccountCache() database.AccountDatabase {
	return &AccountCache{
		pkg_cache.NewCache[string, *account.Account](),
	}
}

func (c *AccountCache) Create(account *account.Account) error {
	ok := c.Cache.Create(account.GetId(), account)
	if !ok {
		return database.ErrAlreadyExists{}
	}
	return nil
}

func (c *AccountCache) Get(id string) (*account.Account, error) {
	account, ok := c.Cache.Get(id)
	if !ok {
		return nil, database.ErrNotExist{}
	}
	return account, nil
}

func (c *AccountCache) Update(account *account.Account) error {
	ok := c.Cache.Update(account.GetId(), account)
	if !ok {
		return database.ErrAlreadyExists{}
	}
	return nil
}

func (c *AccountCache) Delete(id string) error {
	ok := c.Cache.Delete(id)
	if !ok {
		return database.ErrNotExist{}
	}
	return nil
}
