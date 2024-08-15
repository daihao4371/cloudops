package models

import (
	"database/sql/driver"
	"strings"
)

type StringArray []string

// 将字符串类型的参数val按"|"分隔符分割，并返回StringArray类型
func (m *StringArray) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*m = ss
	return nil
}

// 将StringArray类型的m拼接成一个由"|"分隔的字符串，并返回driver.Value类型
func (m StringArray) Value() (driver.Value, error) {
	str := strings.Join(m, "|")
	return str, nil
}
