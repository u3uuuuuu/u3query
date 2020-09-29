package models

import "fmt"

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

