package models

import (
	"errors"
	"sync"
	"time"
	"u3.com/u3query/tree"
)

const (
    SplitLength = 100000
	CacheLength = 3
)


var CacheBt = CacheU{
	cache : make(map[string]*tree.BTree, CacheLength),
	cacheTime: make(map[string]int, CacheLength),
}

//现在只是简单用一个3元素的map作为一个cache， 以后可以根据LRU最少访问原则优化
type CacheU struct {
	mu sync.Mutex
	cache map[string]*tree.BTree
	cacheTime map[string]int
	MaxPrimary int
}

//用于Cache更新
func (c CacheU) Put(cacheKey string, bt *tree.BTree) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Nanosecond()

	if len(c.cache) < CacheLength {
		c.cache[cacheKey] = bt
		c.cacheTime[cacheKey] = now
	} else {

		min, minKey := now, ""
		for k, v := range c.cacheTime {
			if v < min {
				min = v
				minKey = k
			}
		}
		if minKey == ""{
			return errors.New("cache error!!")
		}

		delete(c.cache, minKey)

		c.cache[cacheKey] = bt
		c.cacheTime[cacheKey] = now
	}
	return nil
}

//用于cache获取
func (c CacheU) GetCacheBt(cacheKey string) (*tree.BTree, error) {
	if belongBTree, ok := c.cache[cacheKey]; ok {
		return belongBTree, nil
	} else {
		//如果cache中没有的话，查询持久化文件，如果持久话文件也没有的话，返回空
		rlt, err := tree.ReadBTreeFile(cacheKey)
		if err != nil {
			return nil, err
		}
		c.Put(cacheKey, rlt)
		return rlt, nil
	}
}


func (c CacheU) FlushToDisk(key string, bt *tree.BTree) error {
	_, err := tree.SaveToDisk(bt, key)
	if err != nil {
		return err
	}
	return nil
}

