package user

import (
	"context"

	"github.com/Calmantara/go-user/pkg/domain/response"
)

type UserUsecase interface {
	GetUsersSvc(ctx context.Context, query UserQuery, users *[]*User) (err response.ErrorResponse)
	GetUserByIdSvc(ctx context.Context, user *User) (err response.ErrorResponse)
	UpdateUserSvc(ctx context.Context, user *User) (err response.ErrorResponse)
	InsertUserSvc(ctx context.Context, user *User) (err response.ErrorResponse)
}
