package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Calmantara/go-user/common/logger"
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"github.com/Calmantara/go-user/pkg/usecase/crypto"
)

type UserUsecaseImpl struct {
	sugar    logger.CustomLogger
	userRepo user.UserRepo
	hash     crypto.Hash
}

func NewUserUsecase(sugar logger.CustomLogger, userRepo user.UserRepo, hash crypto.Hash) user.UserUsecase {
	sugar.Logger().Info("init user usecase . . .")

	return &UserUsecaseImpl{sugar: sugar, userRepo: userRepo, hash: hash}
}

func (u *UserUsecaseImpl) GetUsersSvc(ctx context.Context, users []*user.User) (err response.ErrorResponse) {
	return err
}
func (u *UserUsecaseImpl) GetUserByIdSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	return err
}
func (u *UserUsecaseImpl) UpdateUserSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	return err
}
func (u *UserUsecaseImpl) InsertUserSvc(ctx context.Context, userDet *user.User) (err response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T InsertUserSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T InsertUserSvc executed", u)
	// hash password
	hashed, er := u.hash.GeneratePassword(userDet.Password)
	if er != nil {
		u.sugar.WithContext(ctx).Error("error hashing password")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}
	userDet.Password = hashed

	// masking and check credit card and cvv
	if len(userDet.CreditCard.Number) != 12 {
		u.sugar.WithContext(ctx).Error("invalid credit card")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}
	userDet.CreditCard.Number = fmt.Sprintf("********%s", userDet.CreditCard.Number[8:12])
	if len(userDet.CreditCard.Cvv) != 3 {
		u.sugar.WithContext(ctx).Error("invalid cvv number")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}
	userDet.CreditCard.Cvv = fmt.Sprintf("**%s", userDet.CreditCard.Cvv[2:])
	// check expired card
	exp := strings.Split(userDet.CreditCard.Expired, "/")
	if len(exp) != 2 {
		u.sugar.WithContext(ctx).Error("invalid expired number")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}
	month, er := strconv.Atoi(exp[0])
	if er != nil || month < 1 || month > 12 {
		u.sugar.WithContext(ctx).Error("invalid expired number")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}
	if er := u.userRepo.InsertUserRepo(ctx, userDet); er != nil {
		u.sugar.WithContext(ctx).Error("error inserting payload")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}
	return err
}
