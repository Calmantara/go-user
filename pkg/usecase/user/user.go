package user

import (
	"context"

	"github.com/Calmantara/go-user/common/logger"
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/Calmantara/go-user/pkg/domain/user"
)

type UserUsecaseImpl struct {
	sugar    logger.CustomLogger
	userRepo user.UserRepo
}

func NewUserUsecase(sugar logger.CustomLogger, userRepo user.UserRepo) user.UserUsecase {
	sugar.Logger().Info("init user usecase . . .")

	return &UserUsecaseImpl{sugar: sugar, userRepo: userRepo}
}

func (u UserUsecaseImpl) GetUsersSvc(ctx context.Context, users []*user.User) (err response.ErrorResponse) {
	return err
}
func (u UserUsecaseImpl) GetUserByIdSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	return err
}
func (u UserUsecaseImpl) UpdateUserSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	return err
}
func (u UserUsecaseImpl) InsertUserSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	return err
}
