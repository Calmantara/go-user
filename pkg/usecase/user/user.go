package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Calmantara/go-user/lib/logger"
	serviceutil "github.com/Calmantara/go-user/lib/service/util"
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

	// local cache
	userMap     map[uint64]user.User
	notRequired map[string]bool
}

func NewUserUsecase(sugar logger.CustomLogger, userRepo user.UserRepo, hash crypto.Hash, jwt crypto.Jwt, util serviceutil.UtilService) user.UserUsecase {
	sugar.Logger().Info("init user usecase . . .")
	// define map
	userMap := make(map[uint64]user.User)
	notRequired := map[string]bool{
		"name":     true,
		"address":  true,
		"password": true,
	}
	return &UserUsecaseImpl{sugar: sugar, userRepo: userRepo, hash: hash, jwt: jwt, util: util, userMap: userMap, notRequired: notRequired}
}

func (u *UserUsecaseImpl) GetUsersSvc(ctx context.Context, query user.UserQuery, users *[]*user.User) (errMsg response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUserByIdSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUserByIdSvc executed", u)

	// query user lists
	if err := u.userRepo.GetUsersRepo(ctx, query, users); err != nil {
		u.sugar.WithContext(ctx).Error("error getting payload")
		return response.ErrorResponse{
			Error: response.INTERNAL_ERROR_MSG,
			Code:  response.INTERNAL_ERROR_CODE,
		}
	}

	// TODO: decode jwt with worker
	for i := range *users {
		// decode first
		claim := u.jwt.VerifyJWT(ctx, (*users)[i].CreditCardToken.Token)

		var cc creditcard.CreditClaim
		u.util.ObjectMapper(&claim, &cc)

		if cc.Issuer != "go-user" {
			u.sugar.WithContext(ctx).Error("error payload token issuer")
			return response.ErrorResponse{
				Error: response.BAD_REQUEST_MSG,
				Code:  response.BAD_REQUEST_CODE,
			}
		}
		(*users)[i].CreditCard = &cc.CreditCard
		(*users)[i].CreditCard.Cvv = ""
	}

	return errMsg
}
func (u *UserUsecaseImpl) GetUserByIdSvc(ctx context.Context, userDet *user.User) (errMsg response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUserByIdSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUserByIdSvc executed", u)

	userCache, ok := u.userMap[userDet.Id]
	if !ok {
		// get from database
		userCache.Id = userDet.Id
		if err := u.userRepo.GetUserByIdRepo(ctx, &userCache); err != nil {
			u.sugar.WithContext(ctx).Error("error getting payload")
			return response.ErrorResponse{
				Error: response.INTERNAL_ERROR_MSG,
				Code:  response.INTERNAL_ERROR_CODE,
			}
		}
		u.userMap[userCache.Id] = userCache
	}
	*userDet = userCache

	// decode first
	if userDet.Email == "" {
		u.sugar.WithContext(ctx).Error("error user not found")
		return response.ErrorResponse{
			Error: response.ErrorMessage("user not found"),
			Code:  response.ErrorCode(400),
		}
	}
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
	u.sugar.WithContext(ctx).Infof("execute %T InsertUserSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T InsertUserSvc executed", u)
	// get user first
	userTmp := user.User{Id: userDet.Id}
	if errMsg := u.GetUserByIdSvc(ctx, &userTmp); errMsg.Error != "" {
		u.sugar.WithContext(ctx).Errorf("error get user:%+v", errMsg)
		return errMsg
	}
	// detect user changes for not required field
	var userField map[string]any
	u.util.ObjectMapper(&userDet, &userField)

	var userTempMap map[string]any
	u.util.ObjectMapper(&userTmp, &userTempMap)

	for k, v := range userField {
		if u.notRequired[k] {
			if v == "" {
				userField[k] = userTempMap[k]
				userDet.Hased = true
			}
		}
	}
	u.util.ObjectMapper(userField, userDet)
	userDet.PassPhoto = true

	// insert to user database
	if errMsg = u.InsertUserSvc(ctx, userDet); errMsg.Error != "" {
		u.sugar.WithContext(ctx).Errorf("error insert user:%+v", errMsg)
	}
	// delete local cache
	delete(u.userMap, userDet.Id)
	return errMsg
}
func (u *UserUsecaseImpl) InsertUserSvc(ctx context.Context, userDet *user.User) (errMsg response.ErrorResponse) {
	u.sugar.WithContext(ctx).Infof("execute %T InsertUserSvc", u)
	defer u.sugar.WithContext(ctx).Infof("%T InsertUserSvc executed", u)
	// check field
	var userField map[string]any
	u.util.ObjectMapper(&userDet, &userField)
	for k, v := range userField {
		if u.notRequired[k] {
			if v == "" {
				return response.ErrorResponse{
					Error: response.ErrorMessage(
						fmt.Sprintf(string(response.MISSING_FIELD_MSG),
							k)),
					Code: response.MISSING_FIELD_CODE,
				}
			}
			continue
		}
	}
	u.util.ObjectMapper(userField, userDet)
	// passed photo checking
	if !userDet.PassPhoto {
		if userDet.Photos == nil {
			return response.ErrorResponse{
				Error: response.ErrorMessage(
					fmt.Sprintf(string(response.MISSING_FIELD_MSG),
						"photos")),
				Code: response.MISSING_FIELD_CODE,
			}
		}

		if len(userDet.Photos) <= 0 {
			return response.ErrorResponse{
				Error: response.ErrorMessage(
					fmt.Sprintf(string(response.MISSING_FIELD_MSG),
						"photos")),
				Code: response.MISSING_FIELD_CODE,
			}
		}
	}

	// hash password
	if !userDet.Hased {
		hashed, err := u.hash.GeneratePassword(userDet.Password)
		if err != nil {
			u.sugar.WithContext(ctx).Error("error hashing password")
			return response.ErrorResponse{
				Error: response.INTERNAL_ERROR_MSG,
				Code:  response.INTERNAL_ERROR_CODE,
			}
		}
		userDet.Password = hashed
	}

	// check creditcard type
	userDet.CreditCard.Type = strings.ToUpper(userDet.CreditCard.Type)
	if !(creditcard.CreditCardTypeMap[userDet.CreditCard.Type]) {
		u.sugar.WithContext(ctx).Error("invalid credit card type")
		return response.ErrorResponse{
			Error: response.INVALID_CREDIT_CARD_MSG,
			Code:  response.INVALID_CREDIT_CARD_CODE,
		}
	}

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
