package ctl

import (
	"context"
	"errors"
)

type key int

var userKey key

type UserInfo struct {
	Id uint `json:"id"`
}

// 通过 Go 的 context 包来传递用户信息
func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}

func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u) //将用户信息 u 以 userKey 为键存储在 context 中
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userKey).(*UserInfo) //断言为对应类型
	return u, ok
}

func InitUserInfo(ctx context.Context) {
	// TOOD 放缓存，之后的用户信息，走缓存
}
