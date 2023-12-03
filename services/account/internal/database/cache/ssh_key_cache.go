package cache

import (
	account "github.com/de1phin/iam/genproto/services/account/api"
	pkg_cache "github.com/de1phin/iam/pkg/cache"
	"github.com/de1phin/iam/services/account/internal/database"
)

type SshKeyCache struct {
	*pkg_cache.Cache[string, []*account.SshKey]
	fingerprintToAccountId map[string]string
}

func NewSshKeyCache() database.SshKeyDatabase {
	return &SshKeyCache{
		Cache:                  pkg_cache.NewCache[string, []*account.SshKey](),
		fingerprintToAccountId: make(map[string]string),
	}
}

func (c *SshKeyCache) Create(key *account.SshKey) error {
	list, ok := c.Cache.Get(key.GetAccountId())
	if !ok {
		c.Cache.Create(key.GetAccountId(), []*account.SshKey{key})
		c.fingerprintToAccountId[key.GetFingerprint()] = key.GetAccountId()
		return nil
	}

	for _, k := range list {
		if k.GetFingerprint() == key.GetFingerprint() {
			return database.ErrAlreadyExists{}
		}
	}
	list = append(list, key)
	c.Cache.Update(key.GetAccountId(), list)
	c.fingerprintToAccountId[key.GetFingerprint()] = key.GetAccountId()
	return nil
}

func (c *SshKeyCache) Get(fingerprint string) (*account.SshKey, error) {
	accountId, ok := c.fingerprintToAccountId[fingerprint]
	if !ok {
		return nil, database.ErrNotExist{}
	}

	list, _ := c.Cache.Get(accountId)
	for _, key := range list {
		if key.GetFingerprint() == fingerprint {
			return key, nil
		}
	}

	return nil, database.ErrNotExist{}
}

func (c *SshKeyCache) List(accountId string) ([]*account.SshKey, error) {
	list, ok := c.Cache.Get(accountId)
	if !ok {
		return nil, database.ErrNotExist{}
	}
	return list, nil
}

func (c *SshKeyCache) Delete(accountId, fingerprint string) error {
	list, ok := c.Cache.Get(accountId)
	if !ok {
		return database.ErrNotExist{}
	}
	for idx, k := range list {
		if k.GetFingerprint() != fingerprint {
			continue
		}

		if len(list) == 1 {
			c.Cache.Delete(accountId)
		} else {
			c.Cache.Update(accountId, append(list[:idx], list[idx+1:]...))
		}
		delete(c.fingerprintToAccountId, fingerprint)
		return nil
	}

	return database.ErrNotExist{}
}
