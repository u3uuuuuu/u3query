package models

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	Units map[string]*Unit
)

type Unit struct {
	KeySize   int
	Key       string
	ValueSize int
	Value     string
}

func init() {
	Units = make(map[string]*Unit)
	Units["llkjlkj"] = &Unit{7, "llkjlkj", 8, "llkjlkjj"}
}

func (u *Unit) String() string {
	return fmt.Sprintf("Unit KeySize:%d, Key:%s, ValueSize:%d, Value:%s", u.KeySize, u.Key, u.ValueSize, u.Value)
}


func GetUnit(id int) (*Unit, error) {
	begin, end := (id/SplitLength)*SplitLength, (id/SplitLength + 1)*SplitLength
	cacheKey := strconv.Itoa(begin)+"-"+strconv.Itoa(end)
	bt, err := CacheBt.GetCacheBt(cacheKey)
	if err != nil {
		return nil, err
	}
	if unit, ok := bt.Search(id); ok {
		if u, isok := unit.(*Unit); isok {
			return u, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("not find id = %d's data", id))
}


func InsertUnit(u *Unit) (int, error) {
	id := CacheBt.MaxPrimary
	begin, end := (id/SplitLength)*SplitLength, (id/SplitLength + 1)*SplitLength
	cacheKey := strconv.Itoa(begin)+"-"+strconv.Itoa(end)

	bt, err := CacheBt.GetCacheBt(cacheKey)
	if err != nil {
		return 0, err
	}
	bt.Insert(id, u)
	err = CacheBt.Put(cacheKey, bt)
	if err != nil {
		return 0, err
	}
	err = CacheBt.FlushToDisk(cacheKey, bt)
	if err != nil {
		return 0, err
	}
	CacheBt.MaxPrimary++
	return id,nil
}