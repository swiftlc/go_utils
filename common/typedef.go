package common

import (
	. "github.com/swiftlc/go_utils/micro"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Empty struct{}

type Paging struct {
	PageNo   int `json:"page_no"` //从1开始
	PageSize int `json:"page_size"`
}

func (p Paging) Limit() int {
	return p.PageSize
}

func (p Paging) Offset() int {
	return MAX(0, (p.PageNo-1)*p.PageSize)
}

type CommonRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
