package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Calmantara/go-user/common/logger"
	serviceutil "github.com/Calmantara/go-user/common/service/util"
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"github.com/Calmantara/go-user/pkg/usecase/crypto"
)

type UserUsecaseImpl struct {
	sugar    logger.CustomLogger
	userRepo user.UserRepo
	hash     crypto.Hash
	jwt      crypto.Jwt
	util     serviceutil.UtilService
}

func NewUserUsecase(sugar logger.CustomLogger, userRepo user.UserRepo, hash crypto.Hash, jwt crypto.Jwt, util serviceutil.UtilService) user.UserUsecase {
	sugar.Logger().Info("init user usecase . . .")
	return &UserUsecaseImpl{sugar: sugar, userRepo: userRepo, hash: hash, jwt: jwt, util: util}
}

func (u *UserUsecaseImpl) GetUsersSvc(ctx context.Context, users []*user.User) (errMsg response.ErrorResponse) {
	return errMsg
}
func (u *UserUsecaseImpl) GetUserByIdSvc(ctx context.Context, userDet *user.User) (errMsg response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUserByIdSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUserByIdSvc executed", u)

	// get from database
	if err := u.userRepo.GetUserByIdRepo(ctx, userDet); err != nil {
		u.sugar.WithContext(ctx).Error("error getting payload")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}
	// decode first
	claim := u.jwt.VerifyJWT(ctx, userDet.CreditCardToken.Token)

	var cc creditcard.CreditClaim
	u.util.ObjectMapper(&claim, &cc)

	if cc.Issuer != "go-user" {
		u.sugar.WithContext(ctx).Error("error payload token issuer")
		return response.ErrorResponse{
			Error: response.BAD_REQUEST_MSG,
			Code:  response.BAD_REQUEST_CODE,
		}
	}
	userDet.CreditCard = &cc.CreditCard
	userDet.CreditCard.Cvv = ""
	return errMsg
}
func (u *UserUsecaseImpl) UpdateUserSvc(ctx context.Context, userDet *user.User) (errMsg response.ErrorResponse) {
	return errMsg
}
func (u *UserUsecaseImpl) InsertUserSvc(ctx context.Context, userDet *user.User) (errMsg response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T InsertUserSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T InsertUserSvc executed", u)
	// hash password
	hashed, err := u.hash.GeneratePassword(userDet.Password)
	if err != nil {
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
	month, err := strconv.Atoi(exp[0])
	if err != nil || month < 1 || month > 12 {
		u.sugar.WithContext(ctx).Error("invalid expired number")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}

	// create jwt verifiable data
	claim := creditcard.CreditClaim{
		Subject:    userDet.Email,
		IssuedAt:   time.Now().Unix(),
		Issuer:     "go-user",
		Audience:   "code.test.com",
		Type:       "credit_card",
		CreditCard: *userDet.CreditCard,
	}
	creditJwt, err := u.jwt.CreateJWT(ctx, &claim)
	if err != nil {
		u.sugar.WithContext(ctx).Error("invalid create jwt token")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}
	userDet.CreditCardToken = &creditcard.CreditCardToken{
		Token: creditJwt,
	}

	if er := u.userRepo.InsertUserRepo(ctx, userDet); er != nil {
		u.sugar.WithContext(ctx).Error("error inserting payload")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}
	return errMsg
}
