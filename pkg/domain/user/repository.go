package user

import "context"

type UserRepo interface {
	GetUsersRepo(ctx context.Context, userQuery UserQuery, users []*User) (err error)
	GetUserByIdRepo(ctx context.Context, user *User) (err error)
	UpdateUserRepo(ctx context.Context, user *User) (err error)
	InsertUserRepo(ctx context.Context, user *User) (err error)
}
