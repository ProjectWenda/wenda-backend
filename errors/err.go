package errors

import (
	"errors"
)

var (
	ErrDBQueryFail      = errors.New("query to database failed")
	ErrDBAddFail        = errors.New("failed to add object to db")
	ErrInvalidStructure = errors.New("invalid structure, cannot marshal")
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskAddFail    = errors.New("failed to add task")
	ErrTaskUpdateFail = errors.New("failed to update task")
	ErrTaskDeleteFail = errors.New("failed to delete task")
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAddFail       = errors.New("failed to add user")
	ErrUserUpdateFail    = errors.New("failed to update user")
	ErrUserDeleteFail    = errors.New("failed to delete user")
	ErrDiscordIDNotFound = errors.New("discord ID not found")
)

var (
	ErrOrderNotFound   = errors.New("task order not found")
	ErrOrderAddFail    = errors.New("failed to add task order")
	ErrOrderUpdateFail = errors.New("failed to update task order")
	ErrOrderDeleteFail = errors.New("failed to delete task order")
)
