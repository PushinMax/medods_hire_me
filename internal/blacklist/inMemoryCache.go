package blacklist

import (
	"errors"
)
type Cache struct {
	blacklist map[string]string
}

func newCache() *Cache {
	return &Cache{
		blacklist: make(map[string]string),
	}
}

func (c *Cache) Add(jti, ip string) error {
	if _, ok := c.blacklist[jti]; ok {
		return errors.New("token with this identifier already exists")
	}
	c.blacklist[jti] = ip
	return nil
}


func (c *Cache)	Delete(jti string) error {
	if _, ok := c.blacklist[jti]; ok {
		delete(c.blacklist, jti)
		return nil
	}
	return errors.New("such token does not exist")
}

func (c *Cache)	ContainsAndGetIp(jti string) (string, bool) {
	if v, ok := c.blacklist[jti]; ok {
		return v, ok
	}
	return "", false
}