package domain

import "github.com/blog/blog-community/pkg/errors"

var (
	ErrUserNotFound         = errors.New("E404002", "用户不存在")
	ErrEmailAlreadyExist    = errors.New("E409001", "邮箱已被注册")
	ErrUsernameAlreadyExist = errors.New("E409002", "用户名已被占用")
	ErrInvalidPassword      = errors.New("E400003", "密码错误")
)
